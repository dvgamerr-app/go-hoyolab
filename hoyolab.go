package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
)

type DailyAPI struct {
	Domain string `yaml:"domain"`
	Sign   string `yaml:"sign"`
	Info   string `yaml:"info"`
}

type DailyHoyolab struct {
	Browser  string   `yaml:"browser"`
	ActID    string   `yaml:"act_id"`
	Lang     string   `yaml:"lang"`
	Referer  string   `yaml:"referer"`
	Endpoint string   `yaml:"endpoint"`
	API      DailyAPI `yaml:"api"`
}

type Hoyolab struct {
	Client    *resty.Client
	CookieJar []*http.Cookie
	Daily     *DailyHoyolab
}

type ResInfo struct {
	RetCode int32       `json:"retcode"`
	Message string      `json:"message"`
	Data    ResInfoData `json:"data,omitempty"`
}
type ResInfoData struct {
	TotalSignDay int32  `json:"total_sign_day"`
	Today        string `json:"today"`
	IsSign       bool   `json:"is_sign"`
	FirstBind    bool   `json:"first_bind"`
	IsSub        bool   `json:"is_sub"`
	Region       string `json:"region"`
	MonthLastDay bool   `json:"month_last_day"`
}

// {"retcode":0,"message":"OK","data":{"total_sign_day":11,"today":"2022-10-29","is_sign":false,"first_bind":false,"is_sub":true,"region":"os_asia","month_last_day":false}}

var json jsoniter.API = jsoniter.ConfigCompatibleWithStandardLibrary

func (hoyo *Hoyolab) callRequest(params map[string]string) *resty.Request {
	return hoyo.Client.R().
		SetHeaders(hoyo.generateHeaders()).
		SetQueryParams(params)
}

func (hoyo *Hoyolab) generateHeaders() map[string]string {
	uri, _ := url.Parse(hoyo.Daily.Referer)
	return map[string]string{
		"Accept":          "application/json, text/plain, */*",
		"Accept-Language": "en-US,en;q=0.9,th;q=0.8",
		"Cache-Control":   "no-cache",
		"Connection":      "keep-alive",
		"Content-Type":    "application/json;charset=UTF-8",
		"Pragma":          "no-cache",
		"Referer":         fmt.Sprintf("%s?act_id=%s&lang=%s", hoyo.Daily.Referer, hoyo.Daily.ActID, hoyo.Daily.Lang),
		"Origin":          fmt.Sprintf("%s://%s", uri.Scheme, uri.Host),
	}
}

func (hoyo *Hoyolab) SetCookie(rs []*http.Cookie) {
	hoyo.CookieJar = rs
}

func (hoyo *Hoyolab) CheckDaily() error {
	res, err := hoyo.callRequest(map[string]string{
		"lang": hoyo.Daily.Lang,
	}).Options(hoyo.Daily.Endpoint + hoyo.Daily.API.Sign)
	if err != nil {
		return err
	}

	if res.StatusCode() != 204 {
		return fmt.Errorf("response: %s not supported", res.Status())
	}

	return nil
}

func (hoyo *Hoyolab) DailyInfo() (*ResInfo, error) {
	if len(hoyo.CookieJar) == 0 {
		return nil, fmt.Errorf("noting cookie-jar")
	}

	res, err := hoyo.callRequest(map[string]string{
		"lang":   hoyo.Daily.Lang,
		"act_id": hoyo.Daily.ActID,
	}).SetCookies(hoyo.CookieJar).Get(hoyo.Daily.Endpoint + hoyo.Daily.API.Info)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != 200 {
		return nil, fmt.Errorf("ERROR: status: %s", res.Status())
	}

	log.Println("DailyInfo:", res)

	var resInfo ResInfo
	err = json.Unmarshal(res.Body(), &resInfo)
	if err != nil {
		return nil, fmt.Errorf("ERROR: ResInfo{}")
	}
	return &resInfo, nil
}

func (hoyo *Hoyolab) DailySign() (*resty.Response, error) {
	if len(hoyo.CookieJar) == 0 {
		return nil, fmt.Errorf("noting cookie-jar")
	}

	res, err := hoyo.
		callRequest(map[string]string{
			"lang":   hoyo.Daily.Lang,
			"act_id": hoyo.Daily.ActID,
		}).
		SetCookies(hoyo.CookieJar).
		SetBody(map[string]string{"act_id": hoyo.Daily.ActID}).
		Post(hoyo.Daily.Endpoint + hoyo.Daily.API.Sign)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != 200 {
		return nil, fmt.Errorf("ERROR: status: %s", res.Status())
	}

	log.Println("DailySign:", res)

	// var resInfo ResInfo
	// err = json.Unmarshal(res.Body(), &resInfo)
	// if err != nil {
	// 	return nil, fmt.Errorf("ERROR: ResInfo{}")
	// }
	return res, nil
}
