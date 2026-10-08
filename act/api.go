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

type ActUser struct {
	UserInfo ActUserInfo `json:"user_info"`
}

type ActUserInfo struct {
	UID      string `json:"uid"`
	NickName string `json:"nickname"`
}

var dataUserAgent = map[string]string{
	"chrome":  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
	"edge":    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36 Edg/122.0.0.0",
	"firefox": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0",
}

func (e *DailyHoyolab) generateHeaders() map[string]string {
	uri, _ := url.Parse(e.Referer)

	userAgent := dataUserAgent[e.Name]
	if userAgent == "" {
		// Name may already hold a full user-agent string from the profile config.
		userAgent = e.Name
	}
	if userAgent == "" {
		userAgent = dataUserAgent["chrome"]
	}

	headers := map[string]string{
		"Accept":          "application/json, text/plain, */*",
		"Accept-Encoding": "gzip, deflate, br",
		"Accept-Language": "en-US,en;q=0.9,th;q=0.8",
		"Connection":      "keep-alive",
		"Content-Type":    "application/json;charset=UTF-8",
		"Referer":         fmt.Sprintf("%s://%s", uri.Scheme, uri.Host),
		"Origin":          fmt.Sprintf("%s://%s", uri.Scheme, uri.Host),
		"User-Agent":      userAgent,
		"Sec-Fetch-Mode":  `cors`,
		"Sec-Fetch-Site":  `same-site`,
		"Sec-Fetch-Dest":  `empty`,
	}
	if e.SignGame != "" {
		headers["x-rpc-signgame"] = e.SignGame
	}
	return headers
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
	var res struct {
		RetCode int32           `json:"retcode"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(raw.Body(), &res); err != nil {
		return fmt.Errorf("res::%s", "Unmarshal: ActAPI")
	}

	if res.RetCode != 0 {
		return fmt.Errorf(" %s", res.Message)
	}

	if err := json.Unmarshal(res.Data, actData); err != nil {
		return fmt.Errorf("res::Unmarshal: %v", reflect.TypeOf(actData))
	}

	return nil
}

// execute sends the request and decodes the HoYoLAB envelope into T.
func execute[T any](e *DailyHoyolab, req *resty.Request, method, endpoint string) (*T, error) {
	raw, err := req.Execute(method, endpoint)
	if err != nil {
		return nil, fmt.Errorf("api::%+v", err)
	}
	if raw.StatusCode() != 200 {
		return nil, fmt.Errorf("api::%s", raw.Status())
	}

	var res T
	if err := actResponse(raw, &res); err != nil {
		return nil, fmt.Errorf("api::%+v", err)
	}

	if IsDebug {
		log.Printf("%s::%+v\n", e.Label, res)
	}
	return &res, nil
}

func (e *DailyHoyolab) url(path string) string {
	return e.API.Endpoint + path
}

func (e *DailyHoyolab) GetMonthAward(hoyo *Hoyolab) (*ActAward, error) {
	return execute[ActAward](e, hoyo.ActRequest(e), resty.MethodGet, e.url(e.API.Award))
}

func (e *DailyHoyolab) GetAccountUserInfo(hoyo *Hoyolab) (*ActUser, error) {
	// https://bbs-api-os.hoyolab.com/community/painter/wapi/user/full?t=1683734340156
	req := hoyo.ActRequest(e).
		SetQueryParams(map[string]string{"t": fmt.Sprint(time.Now().UnixMilli())}).
		SetBody(map[string]string{"act_id": e.ActID})
	return execute[ActUser](e, req, resty.MethodGet, "https://bbs-api-os.hoyolab.com/community/painter/wapi/user/full")
}

func (e *DailyHoyolab) GetCheckInInfo(hoyo *Hoyolab) (*ActInfo, error) {
	return execute[ActInfo](e, hoyo.ActRequest(e), resty.MethodGet, e.url(e.API.Info))
}

// DailySignIn signs in and reports whether a captcha challenge was triggered.
func (e *DailyHoyolab) DailySignIn(hoyo *Hoyolab) (bool, error) {
	time.Sleep(1 * time.Second)
	req := hoyo.ActRequest(e).SetBody(map[string]string{"act_id": e.ActID})
	res, err := execute[ActSign](e, req, resty.MethodPost, e.url(e.API.Sign))
	if err != nil {
		return false, err
	}

	risk, ok := (*res)["gt_result"].(map[string]any)
	if !ok {
		return false, nil
	}
	code, _ := risk["risk_code"].(float64)
	return code > 0, nil
}
