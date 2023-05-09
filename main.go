package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/zellyn/kooky"
	"github.com/zellyn/kooky/browser/chrome"
)

var configExt string = "yaml"
var logExt string = "log"
var configPath string = ""
var logPath string = ""
var logfile *os.File

func init() {
	filename, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	filename = strings.ReplaceAll(filepath.Base(filename), filepath.Ext(filename), "")

	configPath = fmt.Sprintf("%s.%s", filename, configExt)
	logPath = fmt.Sprintf("%s.%s", filename, logExt)

	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(path.Join(dirname, configExt)); err == nil {
		configPath = path.Join(dirname, configPath)
		logPath = path.Join(dirname, logPath)
	}
	log.SetFlags(log.Lshortfile | log.Ltime)
	if !IsDebug {
		log.SetFlags(log.Ldate | log.Ltime)
		f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(f)
	}
}

func main() {
	if logfile != nil {
		defer logfile.Close()
	}

	// installSchedule := *flag.Bool("install", false, "# install to ")
	// flag.Parse()

	// if installSchedule {
	// 	ex, err := os.Executable()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	exPath := filepath.Dir(ex)
	// 	fmt.Println(exPath)
	// }
	hoyo := GenerateDefaultConfig()
	apiHoyo, err := hoyo.ReadHoyoConfig()
	if err != nil {
		log.Fatal(err)
	}
	hoyo.Daily = apiHoyo

	cookieStore := kooky.FindAllCookieStores()
	log.Printf("Browser total %d sessions", len(cookieStore))
	for _, store := range cookieStore {
		_, err = os.Stat(store.FilePath())
		if os.IsNotExist(err) {
			continue
		}

		for _, act := range hoyo.Browser {
			if len(act.Profile) > 0 && (store.Browser() != act.Browser || !ContainsStrings(act.Profile, store.Profile())) {
				continue
			}
			log.Printf("Allow '%s' :: '%s'", store.Browser(), store.Profile())

			_, err := chrome.CookieJar(store.FilePath())
			if err != nil {
				log.Fatal(err)
			}

		}

		// 	uri, _ := url.Parse(hoyo.Daily.API.Domain)
		// 	hoyo.SetCookie(cookies.Cookies(uri))

		// 	if len(hoyo.CookieJar) == 0 {
		// 		log.Printf("%s::Cookie is empty, please login hoyolab.com.", store.Browser())
		// 		continue
		// 	}

		// 	resInfo, err := hoyo.DailyInfo()
		// 	if err != nil {
		// 		log.Fatalf("Hoyolab::DailyInfo: %v", err)
		// 	}
		// 	if resInfo.RetCode != 0 {
		// 		log.Printf("Hoyolab::DailyInfo: %v", resInfo.Message)
		// 		continue
		// 	}
		// 	// log.Printf("Hoyolab::DailyInfo:%+v", resInfo.Data)

		// 	actInfo, ok := resInfo.Data.(ActInfo)
		// 	if !ok {
		// 		log.Fatalf("DailyInfo: %v", err)
		// 	}
		// 	if actInfo.IsSign {
		// 		log.Printf("Hoyolab::DailyInfo:Claimed Today %s (Total %d Claims)", actInfo.Today, actInfo.TotalSignDay)
		// 		continue
		// 	}

		// 	_, err = hoyo.DailySign()
		// 	if err != nil {
		// 		log.Fatalf("DailySign: %v", err)
		// 	}
		// 	// log.Printf("Hoyolab::DailySign:%+v", resInfo.Data)
	}
}

func ContainsStrings(a []string, x string) bool {
	sort.Strings(a)
	for _, v := range a {
		if v == x {
			return true
		}
	}
	return false
}

func GenerateDefaultConfig() *Hoyolab {
	// Genshin Impact
	//
	// {"code":"ok"}
	// {"retcode":0,"message":"OK","data":{"total_sign_day":11,"today":"2022-10-29","is_sign":false,"first_bind":false,"is_sub":true,"region":"os_asia","month_last_day":false}}
	var apiGenshinImpact = &DailyHoyolab{
		CookieJar: []*http.Cookie{},
		ActID:     "e202102251931481",
		API: DailyAPI{
			Endpoint: "https://sg-hk4e-api.hoyolab.com",
			Domain:   "https://hoyolab.com",
			Award:    "/event/sol/home",
			Info:     "/event/sol/info",
			Sign:     "/event/sol/sign",
		},
		Lang:    "en-us",
		Referer: "https://act.hoyolab.com/ys/event/signin-sea-v3/index.html",
	}

	// Honkai StarRail
	//
	// https://sg-public-api.hoyolab.com/event/luna/os/home?lang=en-us&act_id=e202303301540311
	// https://sg-public-api.hoyolab.com/event/luna/os/sign
	// {"retcode":0,"message":"OK","data":{"code":"","risk_code":0,"gt":"","challenge":"","success":0,"is_risk":false}}
	// https://sg-public-api.hoyolab.com/event/luna/os/info
	// {"retcode":0,"message":"OK","data":{"total_sign_day":7,"today":"2023-05-09","is_sign":false,"is_sub":false,"region":"","sign_cnt_missed":1,"short_sign_day":0}}
	//
	var apiHonkaiStarRail = &DailyHoyolab{
		CookieJar: []*http.Cookie{},
		ActID:     "e202303301540311",
		API: DailyAPI{
			Endpoint: "https://sg-public-api.hoyolab.com/",
			Domain:   "https://hoyolab.com",
			Award:    "/event/luna/os/home",
			Info:     "/event/luna/os/info",
			Sign:     "/event/luna/os/sign",
		},
		Lang:    "en-us",
		Referer: "https://act.hoyolab.com/bbs/event/signin/hkrpg/index.html",
	}

	// Honkai Impact 3
	// https: //act.hoyolab.com/bbs/event/signin-bh3/index.html?act_id=e202110291205111
	var apiHonkaiImpact = &DailyHoyolab{
		CookieJar: []*http.Cookie{},
		ActID:     "e202110291205111",
		API: DailyAPI{
			Endpoint: "https://sg-public-api.hoyolab.com",
			Domain:   "https://hoyolab.com",
			Award:    "/event/mani/home",
			Sign:     "/event/mani/sign",
			Info:     "/event/mani/info",
		},
		Lang:    "en-us",
		Referer: "https://act.hoyolab.com/bbs/event/signin-bh3/index.html",
	}
	return &Hoyolab{
		Client: resty.New(),
		Browser: []BrowserProfile{
			{
				Browser:   "chrome",
				Profile:   []string{"dvgamer", "dvgamerr"},
				UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36",
			},
		},
		Daily: []*DailyHoyolab{apiGenshinImpact, apiHonkaiStarRail, apiHonkaiImpact},
	}
}
