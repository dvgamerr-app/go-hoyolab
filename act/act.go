package act

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/go-resty/resty/v2"
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
	Label     string         `yaml:"label"`
	ActID     string         `yaml:"act_id"`
	API       DailyAPI       `yaml:"api"`
	Lang      string         `yaml:"lang"`
	Referer   string         `yaml:"referer"`
	UserAgent string         `yaml:"user-agent,omitempty"`
	CookieJar []*http.Cookie `yaml:"cookie,omitempty"`
}

type BrowserProfile struct {
	Browser   string   `yaml:"browser"`
	Name      []string `yaml:"name"`
	UserAgent string   `yaml:"userAgent"`
}

type LineNotify struct {
	Token string `yaml:"token"`
	Mini  bool   `yaml:"mini"`
}

type Hoyolab struct {
	Client  *resty.Client    `yaml:"client,omitempty"`
	Notify  LineNotify       `yaml:"notify"`
	Delay   int32            `yaml:"delay"`
	Browser []BrowserProfile `yaml:"profile"`
	Daily   []*DailyHoyolab  `yaml:"config"`
}

const DEBUG string = "DEBUG"

var IsDebug bool = false
var IsDev bool = false

func init() {
	goenv.Load()
	IsDev = os.Getenv(DEBUG) != ""
	IsDebug = os.Getenv(DEBUG) != "" && os.Getenv(DEBUG) != "production"
}

func (hoyo *Hoyolab) WriteHoyoConfig(configPath string) error {
	hoyo.Client = nil
	raw, err := yaml.Marshal(hoyo)
	if err != nil {
		return fmt.Errorf("yaml Marshal::%s", err)
	}
	err = os.WriteFile(configPath, raw, 0644)
	if err != nil {
		return fmt.Errorf("yaml Write::%s", err)
	}
	return nil
}

func (hoyo *Hoyolab) ReadHoyoConfig(configPath string) error {
	raw, err := os.ReadFile(configPath)
	if err != nil {
		log.Println("Default configuration created.")
		hoyo.WriteHoyoConfig(configPath)
		return nil
	}

	var readHoyo *Hoyolab
	err = yaml.Unmarshal(raw, &readHoyo)
	if err != nil {
		log.Println("Default configuration created.")
		hoyo.WriteHoyoConfig(configPath)
		return nil
	} else {
		log.Println("Configuration readed.")
		hoyo.Notify = readHoyo.Notify
		hoyo.Delay = readHoyo.Delay
		hoyo.Browser = readHoyo.Browser
		hoyo.Daily = readHoyo.Daily
		return nil
	}
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
		if jar.Name == "ltoken" {
			return true
		}
	}
	return false
}
