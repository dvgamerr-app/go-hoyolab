package main

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type DailyHoyolab struct {
	ActID    string `yaml:"act_id"`
	Endpoint string `yaml:"endpoint"`
	Sign     string `yaml:"sign"`
	Info     string `yaml:"info"`
	Lang     string `yaml:"lang"`
}

var client *resty.Client = resty.New()

func (hoyo *DailyHoyolab) request(params map[string]string) *resty.Request {
	return client.R().
		SetHeaders(hoyo.generateHeaders()).
		SetQueryParams(params)
}

func (hoyo *DailyHoyolab) generateHeaders() map[string]string {

	return map[string]string{
		"Accept":          "application/json, text/plain, */*",
		"Accept-Language": "en-US,en;q=0.5",
		"Cache-Control":   "max-age=0",
		"Connection":      "keep-alive",
		"Referer":         fmt.Sprintf("https://act.hoyolab.com/ys/event/signin-sea-v3/index.html?act_id=%s&lang=%s", hoyo.ActID, hoyo.Lang),
		"Origin":          "https://webstatic-sea.mihoyo.com",
	}
}

func (hoyo *DailyHoyolab) CheckServer() error {
	res, err := hoyo.request(map[string]string{
		"lang": hoyo.Lang,
	}).Options(hoyo.Endpoint + hoyo.Sign)
	if err != nil {
		return err
	}

	if res.StatusCode() != 204 {
		return fmt.Errorf("response: %s not supported", res.Status())
	}

	return nil
}
