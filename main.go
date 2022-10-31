package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"github.com/zellyn/kooky"
	"github.com/zellyn/kooky/browser/chrome"
	"gopkg.in/yaml.v2"
)

var configFile string = "config.yaml"
var configPath string = ""

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)

	configPath = configFile
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(path.Join(dirname, configFile)); err == nil {
		configPath = path.Join(dirname, configFile)
	}
}

func main() {
	installSchedule := *flag.Bool("install", false, "# install to ")
	flag.Parse()

	if installSchedule {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)
		fmt.Println(exPath)
	}

	hoyo := &Hoyolab{
		Client:    resty.New(),
		CookieJar: []*http.Cookie{},
		Daily: &DailyHoyolab{
			ActID: "e202102251931481",
			API: DailyAPI{
				Endpoint: "https://sg-hk4e-api.hoyolab.com",
				Domain:   "https://hoyolab.com",
				Sign:     "/event/sol/sign",
				Info:     "/event/sol/info",
			},
			Browser:   "chrome",
			Lang:      "en-us",
			Referer:   "https://act.hoyolab.com/ys/event/signin-sea-v3/index.html",
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
		},
	}

	if _, err := os.Stat(configPath); err != nil {
		raw, err := yaml.Marshal(map[string]interface{}{
			"config": hoyo.Daily,
		})
		if err != nil {
			log.Panic("yaml Marshal fail")
		}
		err = os.WriteFile(configPath, raw, 0644)
		if err != nil {
			log.Panic("WriteFile fail")
		}
	}

	for _, store := range kooky.FindAllCookieStores() {
		_, err := os.Stat(store.FilePath())

		if store.Browser() != "chrome" || os.IsNotExist(err) {
			continue
		}

		log.Printf("%s::%s", store.Browser(), store.Profile())
		cookies, err := chrome.CookieJar(store.FilePath())
		if err != nil {
			log.Fatal(err)
		}

		uri, _ := url.Parse(hoyo.Daily.API.Domain)
		hoyo.SetCookie(cookies.Cookies(uri))

		if len(hoyo.CookieJar) == 0 {
			log.Fatalf("%s::Cookie is empty, please login hoyolab.com.", store.Browser())
		}

		resInfo, err := hoyo.DailyInfo()
		if err != nil || resInfo.RetCode != 0 {
			log.Fatalf("Hoyolab::DailyInfo: %v", err)
		}
		// log.Printf("Hoyolab::DailyInfo:%+v", resInfo.Data)

		actInfo, ok := resInfo.Data.(ActInfo)
		if !ok {
			log.Fatalf("DailyInfo: %v", err)
		}
		if actInfo.IsSign {
			log.Printf("Hoyolab::DailyInfo:Claimed Today %s (Total %d Claims)", actInfo.Today, actInfo.TotalSignDay)
			continue
		}

		_, err = hoyo.DailySign()
		if err != nil {
			log.Fatalf("DailySign: %v", err)
		}
		// log.Printf("Hoyolab::DailySign:%+v", resInfo.Data)
	}
}
