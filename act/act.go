package act

import (
	"net/http"
	"net/url"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/tmilewski/goenv"
	"gopkg.in/yaml.v2"
)

type DailyAPI struct {
	Endpoint string `yaml:"endpoint"`
	Domain   string `yaml:"domain"`
	Award    string `yaml:"award"`
	Info     string `yaml:"info"`
	Sign     string `yaml:"sign"`
}

type DailyHoyolab struct {
	Label    string   `yaml:"label"`
	ActID    string   `yaml:"act_id"`
	API      DailyAPI `yaml:"api"`
	Lang     string   `yaml:"lang"`
	Referer  string   `yaml:"referer"`
	SignGame string   `yaml:"signgame,omitempty"`

	// Name is the user-agent key selected from the browser profile at runtime.
	Name      string         `yaml:"-"`
	CookieJar []*http.Cookie `yaml:"-"`
}

type BrowserProfile struct {
	Browser   string   `yaml:"browser"`
	Name      []string `yaml:"name"`
	UserAgent string   `yaml:"userAgent"`
}

type LineNotify struct {
	Token      string `yaml:"token,omitempty"` // legacy alias of line-notify
	LINENotify string `yaml:"line-notify"`
	Discord    string `yaml:"discord-webhook"`
	Mini       bool   `yaml:"mini"`
}

type Hoyolab struct {
	Client  *resty.Client    `yaml:"-"`
	Notify  LineNotify       `yaml:"notify"`
	Delay   int32            `yaml:"delay"`
	Profile []BrowserProfile `yaml:"profile"`
	Games   []string         `yaml:"games"`
	// Legacy is the pre-`games` format; it is read but never written back.
	Legacy []*DailyHoyolab `yaml:"config,omitempty"`
	Daily  []*DailyHoyolab `yaml:"-"`
}

const DEBUG string = "DEBUG"

var IsDebug bool = false
var IsDev bool = false

func init() {
	if err := goenv.Load(); err != nil {
		log.Error().Err(err).Msg("load .env")
	}
	IsDev = os.Getenv(DEBUG) != ""
	IsDebug = os.Getenv(DEBUG) != "" && os.Getenv(DEBUG) != "production"
}

// resolveDaily builds Daily from Games/Legacy and records the labels in Games
// so legacy configurations are migrated to the compact format on write.
func (hoyo *Hoyolab) resolveDaily() {
	games, unknown := ResolveGames(hoyo.Games, hoyo.Legacy)
	for _, label := range unknown {
		log.Warn().Msgf("Unknown game '%s' ignored.", label)
	}
	hoyo.Daily = games
	hoyo.Games = hoyo.Games[:0:0]
	for _, game := range games {
		hoyo.Games = append(hoyo.Games, game.Label)
	}
}

func (hoyo *Hoyolab) WriteHoyoConfig(configPath string) {
	hoyo.Legacy = nil
	raw, err := yaml.Marshal(hoyo)
	if err != nil {
		log.Fatal().Msgf("yaml Marshal::%s", err)
	}
	if err = os.WriteFile(configPath, raw, 0644); err != nil {
		log.Fatal().Msgf("yaml Write::%s", err)
	}
}

func (hoyo *Hoyolab) ReadHoyoConfig(configPath string) error {
	raw, err := os.ReadFile(configPath)
	if err != nil {
		log.Info().Msg("Default configuration created.")
		hoyo.resolveDaily()
		hoyo.WriteHoyoConfig(configPath)
		return nil
	}

	var read Hoyolab
	if err = yaml.Unmarshal(raw, &read); err != nil {
		return err
	}

	log.Info().Msg("Readed configuration .")
	hoyo.Notify = read.Notify
	if read.Notify.Token != "" {
		hoyo.Notify.LINENotify = read.Notify.Token
		hoyo.Notify.Token = ""
	}
	hoyo.Delay = read.Delay
	hoyo.Profile = read.Profile
	hoyo.Games = read.Games
	hoyo.Legacy = read.Legacy
	hoyo.resolveDaily()
	hoyo.WriteHoyoConfig(configPath)
	return nil
}

func (hoyo *Hoyolab) IsCookieToken(cookies http.CookieJar) bool {
	for _, act := range hoyo.Daily {
		uri, _ := url.Parse(act.API.Domain)
		act.SetCookie(cookies.Cookies(uri))
		if !act.IsCookieLogin() {
			act.CookieJar = nil
			return false
		}
	}
	return true
}

func (e *DailyHoyolab) SetCookie(rs []*http.Cookie) {
	e.CookieJar = rs
}

func (e *DailyHoyolab) IsCookieLogin() bool {
	for _, jar := range e.CookieJar {
		if jar.Name == "ltoken" || jar.Name == "ltoken_v2" {
			return true
		}
	}
	return false
}
