package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"hoyolab/act"

	"github.com/alexflint/go-arg"
	"github.com/dvgamerr-app/go-kooky/v2/browser"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const configExt = "yaml"

var (
	configPath string
	cli        struct {
		Run         bool   `arg:"--run,-r" default:"true"`
		ListBrowser bool   `arg:"--list,-l"`
		Token       string `arg:"--token" help:"set token one job"`
		Generate    string `arg:"--gen" help:"want value 'browser:profile'"`
	}
)

func init() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	execFilename, err := os.Executable()
	if err != nil {
		log.Fatal().Err(err).Msg("resolve executable")
	}
	baseFilename := strings.TrimSuffix(filepath.Base(execFilename), filepath.Ext(execFilename))
	dirname := filepath.Dir(execFilename)

	timeFormat := time.RFC3339
	if act.IsDev {
		timeFormat = time.TimeOnly
		if dirname, err = os.Getwd(); err != nil {
			log.Fatal().Err(err).Msg("resolve working directory")
		}
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: timeFormat})

	configPath = filepath.Join(dirname, fmt.Sprintf("%s.%s", baseFilename, configExt))
}

func main() {
	arg.MustParse(&cli)

	hoyo := newDefaultConfig()
	if err := hoyo.ReadHoyoConfig(configPath); err != nil {
		log.Fatal().Err(err).Msg("read configuration")
	}

	for _, profile := range hoyo.Profile {
		if unsupportedBrowsers[profile.Browser] {
			log.Warn().Msgf("profile browser '%s' is no longer supported (cookies use App-Bound Encryption); set it to 'firefox' and log in to HoYoLAB there.",
				profile.Browser)
		}
	}

	cookieStore := findHoyolabStores(hoyo)
	log.Info().Msgf("Browser total %d sessions", len(cookieStore))
	if len(cookieStore) == 0 {
		log.Fatal().Msgf("Not found HoYoLAB login in any supported browser (%s). Chrome, Edge and Brave are not supported; log in to HoYoLAB with Firefox.",
			strings.Join(readableBrowsers, ", "))
	}

	if cli.ListBrowser {
		log.Info().Msg("List of profiles for config:")
		for _, store := range cookieStore {
			key := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", store.Browser, store.Profile)))
			log.Info().Msgf(" - %s:%s [%s]", store.Browser, store.Profile, key)
		}
		return
	}

	var notifyMessage []string
	for _, store := range cookieStore {
		for _, profile := range hoyo.Profile {
			if profile.Browser != "" && store.Browser != profile.Browser {
				continue
			}
			if len(profile.Name) > 0 && !slices.Contains(profile.Name, store.Profile) {
				continue
			}
			log.Info().Msgf(`Profile %s:"%s"`, store.Browser, store.Profile)

			if msg := runProfile(hoyo, store, profile); msg != "" {
				if len(hoyo.Profile) > 1 {
					notifyMessage = append(notifyMessage, "\n")
				}
				notifyMessage = append(notifyMessage, msg)
			}
		}
	}

	if err := hoyo.NotifyMessage(strings.Join(notifyMessage, "\n")); err != nil {
		log.Info().Msgf("NotifyMessage  : %v", err)
	}
}

// readableBrowsers lists the browsers whose cookies are read. Firefox-family
// browsers store cookies unencrypted. Opera is Chromium-based and is attempted
// on a best-effort basis; a v20 store is reported and skipped.
var readableBrowsers = []string{"firefox", "librewolf", "waterfox", "zen", "floorp", "opera"}

// unsupportedBrowsers are Chromium browsers whose cookies use App-Bound
// Encryption (v20) and can no longer be read.
var unsupportedBrowsers = map[string]bool{"chrome": true, "edge": true, "brave": true}

// browserStore is one browser profile that holds a HoYoLAB login.
type browserStore struct {
	Browser string
	Profile string
	jar     *storeJar
}

// storeJar adapts a go-kooky cookie store to http.CookieJar so the existing
// act.IsCookieToken flow keeps working. Reads are cached per host.
type storeJar struct {
	store *browser.CookieStore
	cache map[string][]*http.Cookie
	err   error
}

func newStoreJar(name, profile string) *storeJar {
	return &storeJar{store: browser.NewCookieStore(name, profile), cache: map[string][]*http.Cookie{}}
}

func (j *storeJar) SetCookies(*url.URL, []*http.Cookie) {}

func (j *storeJar) Cookies(u *url.URL) []*http.Cookie {
	host := u.Hostname()
	if cached, ok := j.cache[host]; ok {
		return cached
	}
	cookies, err := j.store.CookiesForHost(host)
	if err != nil {
		j.err = err
	}
	j.cache[host] = cookies
	return cookies
}

// findHoyolabStores returns the browser profiles that hold a HoYoLAB login.
func findHoyolabStores(hoyo *act.Hoyolab) []browserStore {
	var found []browserStore
	for _, name := range readableBrowsers {
		profiles, err := browser.GetProfiles(name)
		if err != nil {
			log.Debug().Msgf("%s: not available: %v", name, err)
			continue
		}
		for _, profile := range profiles {
			jar := newStoreJar(name, profile)
			if hoyo.IsCookieToken(jar) {
				found = append(found, browserStore{Browser: name, Profile: profile, jar: jar})
			} else if errors.Is(jar.err, browser.ErrAppBoundEncryption) {
				log.Warn().Msgf(`%s:"%s" skipped: %v`, name, profile, jar.err)
			}
		}
	}
	return found
}

// runProfile checks in every configured game for one browser profile and
// returns the notification text, or "" when nothing was claimed.
func runProfile(hoyo *act.Hoyolab, store browserStore, profile act.BrowserProfile) string {
	// Cookies are shared on hoyo.Daily, so reload them for this store.
	if !hoyo.IsCookieToken(store.jar) {
		return ""
	}
	hoyo.Client = resty.New()

	var user *act.ActUser
	var daySign int32 = -1
	var awards []string

	for _, game := range hoyo.Daily {
		game.Name = profile.UserAgent

		if user == nil {
			res, err := game.GetAccountUserInfo(hoyo)
			if err != nil {
				log.Info().Msgf("%s::GetUserInfo    : %v", game.Label, err)
				continue
			}
			user = res
			log.Info().Msgf("%s::GetUserInfo    : Hi, '%s'", game.Label, user.UserInfo.NickName)
		}

		resAward, err := game.GetMonthAward(hoyo)
		if err != nil {
			log.Info().Msgf("%s::GetMonthAward  : %v", game.Label, err)
			continue
		}

		resInfo, err := game.GetCheckInInfo(hoyo)
		if err != nil {
			log.Info().Msgf("%s::GetCheckInInfo :%v", game.Label, err)
			continue
		}

		log.Info().Msgf("%s::GetCheckInInfo : Checked in for %d days", game.Label, resInfo.TotalSignDay)
		if resInfo.IsSign {
			log.Info().Msgf("%s::DailySignIn    : Claimed %s", game.Label, resInfo.Today)
			continue
		}

		isRisk, err := game.DailySignIn(hoyo)
		if err != nil {
			log.Info().Msgf("%s::DailySignIn    : %v", game.Label, err)
			continue
		}
		if daySign < 0 {
			daySign = resInfo.TotalSignDay + 1
		}

		day := int(resInfo.TotalSignDay)
		if day+1 >= len(resAward.Awards) {
			log.Info().Msgf("%s::GetMonthAward  : No award for day %d", game.Label, day+1)
			continue
		}
		award := resAward.Awards[day+1]
		log.Info().Msgf("%s::GetMonthAward  : Today's received %s x%d", game.Label, award.Name, award.Count)

		awards = append(awards, formatAward(hoyo.Notify.Mini, game.Label, resInfo.TotalSignDay+1, isRisk, award))
	}

	if len(awards) == 0 || user == nil {
		return ""
	}
	if hoyo.Notify.Mini {
		return fmt.Sprintf("%s, at day %d your got %s", user.UserInfo.NickName, daySign, strings.Join(awards, ", "))
	}
	return fmt.Sprintf("\nHi, %s Checked in for %d days.\n%s", user.UserInfo.NickName, daySign, strings.Join(awards, "\n"))
}

func formatAward(mini bool, label string, day int32, isRisk bool, award act.ActAwardItem) string {
	switch {
	case mini && isRisk:
		return fmt.Sprintf("Challenge captcha (%s)", label)
	case mini:
		return fmt.Sprintf("*%s x%d* (%s)", award.Name, award.Count, label)
	case isRisk:
		return fmt.Sprintf("*[%s]* at day %d challenge captcha", label, day)
	default:
		return fmt.Sprintf("*[%s]* at day %d received %s x%d", label, day, award.Name, award.Count)
	}
}

func newDefaultConfig() *act.Hoyolab {
	return &act.Hoyolab{
		Notify: act.LineNotify{Mini: true},
		Delay:  150,
		Profile: []act.BrowserProfile{
			{
				Browser:   "firefox",
				Name:      []string{},
				UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0",
			},
		},
	}
}
