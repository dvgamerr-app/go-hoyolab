package act

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"time"

	"github.com/go-resty/resty/v2"
)

type ActAPI struct {
	RetCode int32  `json:"retcode"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ActAward struct {
	Month      int16          `json:"month"`
	Awards     []ActAwardItem `json:"awards"`
	Biz        string         `json:"biz"`
	Resign     bool           `json:"resign"`
	ExtraAward struct {
		HasExtraAward  bool   `json:"has_extra_award"`
		StartTime      string `json:"start_time"`
		EndTime        string `json:"end_time"`
		List           []any  `json:"list"`
		StartTimestamp string `json:"start_timestamp"`
		EndTimestamp   string `json:"end_timestamp"`
	} `json:"short_extra_award"`
}

type ActAwardItem struct {
	Icon  string `json:"icon"`
	Name  string `json:"name"`
	Count int32  `json:"cnt"`
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

type ActSign map[string]any

// var json jsoniter.API = jsoniter.ConfigCompatibleWithStandardLibrary

func (e *DailyHoyolab) generateHeaders() map[string]string {
	uri, _ := url.Parse(e.Referer)
	return map[string]string{
		"Accept":          "application/json, text/plain, */*",
		"Accept-Language": "en-US,en;q=0.9,th;q=0.8",
		"Connection":      "keep-alive",
		"Content-Type":    "application/json;charset=UTF-8",
		"Referer":         fmt.Sprintf("%s://%s", uri.Scheme, uri.Host),
		"Origin":          fmt.Sprintf("%s://%s", uri.Scheme, uri.Host),
		"User-Agent":      e.UserAgent,
	}
}

func (hoyo *Hoyolab) ActRequest(act *DailyHoyolab) *resty.Request {
	// delay request 150ms default.
	time.Sleep(time.Duration(hoyo.Delay) * time.Millisecond)

	return hoyo.Client.R().
		SetCookies(act.CookieJar).
		SetHeaders(act.generateHeaders()).
		SetQueryParams(map[string]string{"lang": act.Lang, "act_id": act.ActID})
}

func actResponse[T any](raw *resty.Response, actData *T) error {
	var res ActAPI
	if err := json.Unmarshal(raw.Body(), &res); err != nil {
		return fmt.Errorf("res::%s", "Unmarshal: ActAPI")
	}

	if res.RetCode != 0 {
		return fmt.Errorf("res::RetCode: %s", res.Message)
	}

	data, err := json.Marshal(res.Data)
	if err != nil {
		return fmt.Errorf("res::Marshal: %v", reflect.TypeOf(actData))
	}

	if err := json.Unmarshal(data, actData); err != nil {
		return fmt.Errorf("res::Unmarshal: %v", reflect.TypeOf(actData))
	}

	return nil
}

func (e *DailyHoyolab) GetMonthAward(hoyo *Hoyolab) (*ActAward, error) {
	raw, err := hoyo.ActRequest(e).Get(fmt.Sprintf("%s%s", e.API.Endpoint, e.API.Award))
	if err != nil {
		return nil, fmt.Errorf("hoyo::%+v", err)
	}
	if raw.StatusCode() != 200 {
		return nil, fmt.Errorf("hoyo::%s", raw.Status())
	}

	var res ActAward
	if err := actResponse(raw, &res); err != nil {
		return nil, fmt.Errorf("hoyo::%+v", err)
	}

	if IsDebug {
		log.Printf("%s::%+v\n", e.Label, map[string]any{
			"Month":      res.Month,
			"Biz":        res.Biz,
			"Resign":     res.Resign,
			"ExtraAward": res.ExtraAward,
		})
	}
	return &res, nil
}

func (e *DailyHoyolab) GetCheckInInfo(hoyo *Hoyolab) (*ActInfo, error) {
	raw, err := hoyo.ActRequest(e).Get(fmt.Sprintf("%s%s", e.API.Endpoint, e.API.Info))
	if err != nil {
		return nil, fmt.Errorf("hoyo::%+v", err)
	}
	if raw.StatusCode() != 200 {
		return nil, fmt.Errorf("hoyo::%s", raw.Status())
	}

	var res ActInfo
	if err := actResponse(raw, &res); err != nil {
		return nil, fmt.Errorf("hoyo::%+v", err)
	}

	if IsDebug {
		log.Printf("%s::%+v\n", e.Label, res)
	}
	return &res, nil
}

func (e *DailyHoyolab) DailySignIn(hoyo *Hoyolab) (*ActSign, error) {
	raw, err := hoyo.ActRequest(e).SetBody(map[string]string{"act_id": e.ActID}).Post(fmt.Sprintf("%s%s", e.API.Endpoint, e.API.Sign))
	if err != nil {
		return nil, fmt.Errorf("hoyo::%+v", err)
	}
	if raw.StatusCode() != 200 {
		return nil, fmt.Errorf("hoyo::%s", raw.Status())
	}

	var res ActSign
	if err := actResponse(raw, &res); err != nil {
		return nil, fmt.Errorf("hoyo::%+v", err)
	}

	if IsDebug {
		log.Printf("%s::%+v\n", e.Label, res)
	}
	return &res, nil
}
