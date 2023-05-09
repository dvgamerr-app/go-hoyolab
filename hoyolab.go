package main

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

const DEBUG string = "DEBUG"

var IsDebug bool = false

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
	CookieJar []*http.Cookie `yaml:"cookie,omitempty"`
}

type BrowserProfile struct {
	Browser   string   `yaml:"browser"`
	Name      []string `yaml:"name"`
	UserAgent string   `yaml:"userAgent"`
}

type Hoyolab struct {
	Client  *resty.Client    `yaml:"client,omitempty"`
	Browser []BrowserProfile `yaml:"profile"`
	Daily   []*DailyHoyolab  `yaml:"config"`
}

type ActAPI struct {
	RetCode int32       `json:"retcode"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ActInfo struct {
	TotalSignDay int32  `json:"total_sign_day"`
	Today        string `json:"today"`
	IsSign       bool   `json:"is_sign"`
	FirstBind    bool   `json:"first_bind"`
	IsSub        bool   `json:"is_sub"`
	Region       string `json:"region"`
	MonthLastDay bool   `json:"month_last_day"`
}

type ActSign struct {
	Code string `json:"code"`
}

func init() {
	goenv.Load()
	IsDebug = os.Getenv(DEBUG) != "" && os.Getenv(DEBUG) != "production"
}

func (hoyo *Hoyolab) WriteHoyoConfig() error {
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

func (hoyo *Hoyolab) ReadHoyoConfig() error {
	raw, err := os.ReadFile(configPath)
	if err != nil {
		log.Println("Default configuration created.")
		hoyo.WriteHoyoConfig()
		return nil
	}

	var readHoyo *Hoyolab
	err = yaml.Unmarshal(raw, &readHoyo)
	if err != nil {
		log.Println("Default configuration created.")
		hoyo.WriteHoyoConfig()
		return nil
	} else {
		log.Println("Configuration readed.")
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

// var json jsoniter.API = jsoniter.ConfigCompatibleWithStandardLibrary

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

// func (hoyo *DailyHoyolab) ActRequest() *resty.Request {
// 	return hoyo.Client.R().
// 		SetHeaders(hoyo.generateHeaders()).
// 		SetQueryParams(map[string]string{
// 			"lang":   hoyo.Daily.Lang,
// 			"act_id": hoyo.Daily.ActID,
// 		})
// }

// func (hoyo *DailyHoyolab) DailyInfo() (*ActAPI, error) {
// 	if len(hoyo.CookieJar) == 0 {
// 		return nil, fmt.Errorf("hoyo::%s", "DailyInfo - CookieJar is empty")
// 	}

// 	raw, err := hoyo.ActRequest().
// 		SetCookies(hoyo.CookieJar).
// 		Get(hoyo.Daily.API.Endpoint + hoyo.Daily.API.Info)

// 	if err != nil {
// 		return nil, err
// 	}

// 	if raw.StatusCode() != 200 {
// 		return nil, fmt.Errorf("hoyo::%s - %s", "DailyInfo", raw.Status())
// 	}

// 	// log.Println("DailyInfo:\n", raw)

// 	var (
// 		res ActAPI
// 		act ActInfo
// 	)

// 	err = json.Unmarshal(raw.Body(), &res)
// 	if err != nil {
// 		return nil, fmt.Errorf("json::%s", "Unmarshal: ActAPI")
// 	}

// 	data, err := json.Marshal(res.Data)
// 	if err != nil {
// 		return nil, fmt.Errorf("json::%s", "Marshal: interface{}")
// 	}

// 	err = json.Unmarshal(data, &act)
// 	if err != nil {
// 		return nil, fmt.Errorf("json::%s", "Unmarshal: ActAPI")
// 	}

// 	res.Data = act

// 	return &res, nil
// }

// func (hoyo *DailyHoyolab) DailySign() (*ActAPI, error) {
// 	if len(hoyo.CookieJar) == 0 {
// 		return nil, fmt.Errorf("hoyo::%s", "DailyInfo - CookieJar is empty")
// 	}

// 	// if IsDebug {
// 	// 	return &ActAPI{
// 	// 		RetCode: 0,
// 	// 		Message: "OK",
// 	// 		Data:    ActSign{Code: "ok"},
// 	// 	}, nil
// 	// }

// 	raw, err := hoyo.
// 		ActRequest().
// 		SetCookies(hoyo.CookieJar).
// 		SetBody(map[string]string{"act_id": hoyo.Daily.ActID}).
// 		Post(hoyo.Daily.API.Endpoint + hoyo.Daily.API.Sign)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if raw.StatusCode() != 200 {
// 		return nil, fmt.Errorf("hoyo::%s - %s", "DailyInfo", raw.Status())
// 	}

// 	// log.Println("DailySign:\n", raw)

// 	var resSign ActAPI
// 	err = json.Unmarshal(raw.Body(), &resSign)
// 	if err != nil {
// 		return nil, fmt.Errorf("ERROR: ResInfo{}")
// 	}
// 	return &resSign, nil
// }

// func (hoyo *DailyHoyolab) generateHeaders() map[string]string {
// 	uri, _ := url.Parse(hoyo.Daily.Referer)
// 	return map[string]string{
// 		"Accept":          "application/json, text/plain, */*",
// 		"Accept-Language": "en-US,en;q=0.9,th;q=0.8",
// 		"Cache-Control":   "no-cache",
// 		"Connection":      "keep-alive",
// 		"Content-Type":    "application/json;charset=UTF-8",
// 		"Pragma":          "no-cache",
// 		"Referer":         fmt.Sprintf("%s?act_id=%s&lang=%s", hoyo.Daily.Referer, hoyo.Daily.ActID, hoyo.Daily.Lang),
// 		"Origin":          fmt.Sprintf("%s://%s", uri.Scheme, uri.Host),
// 		"User-Agent":      hoyo.Daily.UserAgent,
// 		"referrerPolicy":  "strict-origin-when-cross-origin",
// 		"mode":            "cors",
// 		"credentials":     "include",
// 	}
// }
