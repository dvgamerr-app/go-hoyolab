package main

import (
	"fmt"
	"hoyolab/act"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestCheckInAccountOptionalProfile(t *testing.T) {
	for _, tc := range []struct {
		name     string
		status   int
		profile  string
		nickname string
	}{
		{"forbidden", 403, `Forbidden`, "Browser profile"},
		{"retcode", 200, `{"retcode":-100,"message":"login expired"}`, "Browser profile"},
		{"malformed", 200, `not JSON`, "Browser profile"},
		{"missing profile", 200, `{"retcode":0,"data":null}`, "Browser profile"},
		{"nickname", 200, `{"retcode":0,"data":{"user_info":{"nickname":"Traveler"}}}`, "Traveler"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, mini := range []bool{true, false} {
				hoyo := newTestConfig()
				hoyo.Delay = 0
				hoyo.Notify.Mini = mini
				profiles, signs := 0, 0
				hoyo.Client = resty.New().SetTransport(roundTripFunc(func(r *http.Request) (*http.Response, error) {
					status, body := 200, `{"retcode":0,"data":{}}`
					if r.URL.Host == "bbs-api-os.hoyolab.com" {
						profiles++
						status = tc.status
						body = tc.profile
					} else if strings.HasSuffix(r.URL.Path, "/home") {
						body = `{"retcode":0,"data":{"awards":[{"name":"Last reward","cnt":2}]}}`
					} else if strings.HasSuffix(r.URL.Path, "/info") {
						body = `{"retcode":0,"data":{"total_sign_day":0,"is_sign":false}}`
					} else if strings.HasSuffix(r.URL.Path, "/sign") {
						signs++
						if r.Method != "POST" {
							t.Errorf("sign method %s", r.Method)
						}
					} else {
						t.Fatalf("unexpected request %s", r.URL)
					}
					return &http.Response{StatusCode: status, Status: fmt.Sprintf("%d %s", status, http.StatusText(status)), Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
				}))
				message := checkInAccount(hoyo, "test-agent", "Browser profile")
				if profiles != 1 || signs != 4 {
					t.Fatalf("profile requests=%d, signs=%d", profiles, signs)
				}
				if !strings.Contains(message, tc.nickname) || strings.Count(message, "Last reward") != 4 {
					t.Fatalf("unexpected notification: %q", message)
				}
			}
		})
	}
}

func TestCheckInAccountDailyFailureStillStopsSignIn(t *testing.T) {
	for _, tc := range []struct {
		name, info, sign string
		status           int
		wantSigns        int
	}{
		{"expired", `{"retcode":-100,"message":"Please log in"}`, "", 200, 0},
		{"forbidden", "Forbidden", "", 403, 0},
		{"already claimed", `{"retcode":0,"data":{"is_sign":true}}`, "", 200, 0},
		{"sign rejected", `{"retcode":0,"data":{"total_sign_day":0}}`, `{"retcode":-100,"message":"Please log in"}`, 200, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			hoyo := newTestConfig()
			hoyo.Delay = 0
			hoyo.Daily = hoyo.Daily[:1]
			signs := 0
			hoyo.Client = resty.New().SetTransport(roundTripFunc(func(r *http.Request) (*http.Response, error) {
				status, body := 200, `{"retcode":0,"data":{"awards":[{"name":"reward","cnt":1}]}}`
				switch {
				case r.URL.Host == "bbs-api-os.hoyolab.com":
					status = 403
					body = "Forbidden"
				case strings.HasSuffix(r.URL.Path, "/info"):
					status = tc.status
					body = tc.info
				case strings.HasSuffix(r.URL.Path, "/sign"):
					signs++
					body = tc.sign
				}
				return &http.Response{StatusCode: status, Status: http.StatusText(status), Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
			}))
			if got := checkInAccount(hoyo, "test-agent", ""); got != "" {
				t.Errorf("false success: %q", got)
			}
			if signs != tc.wantSigns {
				t.Errorf("signs=%d, want %d", signs, tc.wantSigns)
			}
		})
	}
}

func TestCheckInAccountProfileTransportFailureAndCaptcha(t *testing.T) {
	hoyo := newTestConfig()
	hoyo.Delay = 0
	hoyo.Daily = hoyo.Daily[:1]
	hoyo.Client = resty.New().SetTransport(roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Host == "bbs-api-os.hoyolab.com" {
			return nil, fmt.Errorf("mock profile unavailable")
		}
		body := `{"retcode":0,"data":{}}`
		if strings.HasSuffix(r.URL.Path, "/sign") {
			body = `{"retcode":0,"data":{"gt_result":{"risk_code":1}}}`
		}
		return &http.Response{StatusCode: 200, Status: "200 OK", Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
	}))
	message := checkInAccount(hoyo, "test-agent", "")
	if !strings.Contains(message, "HoYoLAB user") || !strings.Contains(message, "Challenge captcha") {
		t.Fatalf("notification=%q", message)
	}
	if strings.Contains(message, "received") {
		t.Fatalf("false award success: %q", message)
	}
}
func newTestConfig() *act.Hoyolab {
	hoyo := newDefaultConfig()
	hoyo.Daily, _ = act.ResolveGames([]string{"GSI", "HSR", "ZZZ", "HI3"}, nil)
	return hoyo
}
