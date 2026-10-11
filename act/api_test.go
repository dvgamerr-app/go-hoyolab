package act

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestProfileAndDailyRequestsAreSeparate(t *testing.T) {
	daily := &DailyHoyolab{ActID: "example", Lang: "en-us", Name: "test-agent", Referer: "https://act.hoyolab.com/example", CookieJar: []*http.Cookie{{Name: "ltoken_v2", Value: "test-only"}}, API: DailyAPI{Endpoint: "https://example.invalid", Info: "/info"}}
	daily.SignGame = "zzz"
	requests := 0
	hoyo := &Hoyolab{Client: resty.New().SetTransport(roundTripFunc(func(r *http.Request) (*http.Response, error) {
		requests++
		if r.Header.Get("Cookie") != "ltoken_v2=test-only" {
			t.Errorf("missing test cookie")
		}
		if r.Header.Get("User-Agent") != "test-agent" {
			t.Errorf("missing user agent")
		}
		if r.URL.Host == "bbs-api-os.hoyolab.com" {
			if r.Header.Get("x-rpc-signgame") != "" {
				t.Error("daily signgame header leaked into profile request")
			}
			if r.URL.RawQuery != "" {
				t.Errorf("profile query leaked: %s", r.URL.RawQuery)
			}
			if r.Body != nil {
				b, _ := io.ReadAll(r.Body)
				if len(b) != 0 {
					t.Errorf("profile body leaked: %s", b)
				}
			}
			if r.Header.Get("Referer") != "https://www.hoyolab.com/" {
				t.Errorf("wrong community referer")
			}
		} else {
			if r.Header.Get("x-rpc-signgame") != "zzz" {
				t.Error("daily signgame header was lost")
			}
			if r.URL.Query().Get("act_id") != "example" || r.URL.Query().Get("lang") != "en-us" {
				t.Errorf("missing daily query")
			}
			if r.Header.Get("Referer") != "https://act.hoyolab.com" {
				t.Errorf("wrong daily referer")
			}
		}
		return &http.Response{StatusCode: 200, Status: "200 OK", Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"retcode":0,"data":{}}`)), Request: r}, nil
	}))}
	if _, err := daily.GetAccountUserInfo(hoyo); err != nil {
		t.Fatal(err)
	}
	if _, err := daily.GetCheckInInfo(hoyo); err != nil {
		t.Fatal(err)
	}
	if requests != 2 {
		t.Fatalf("requests=%d", requests)
	}
}

func TestResponseDiagnostics(t *testing.T) {
	for _, tc := range []struct {
		status     int
		body, want string
	}{
		{403, "Forbidden", "403 Forbidden"},
		{200, `{"retcode":-100,"message":"Please log in"}`, "retcode=-100: Please log in"},
		{200, "not JSON", "HTTP 200 OK"},
	} {
		t.Run(tc.want, func(t *testing.T) {
			hoyo := &Hoyolab{Client: resty.New().SetTransport(roundTripFunc(func(r *http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: tc.status, Status: fmt.Sprintf("%d %s", tc.status, http.StatusText(tc.status)), Header: make(http.Header), Body: io.NopCloser(strings.NewReader(tc.body)), Request: r}, nil
			}))}
			_, err := (&DailyHoyolab{}).GetAccountUserInfo(hoyo)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error=%v, want %q", err, tc.want)
			}
		})
	}
}

func TestProfileTimeoutDoesNotChangeDailyClient(t *testing.T) {
	hoyo := &Hoyolab{Client: resty.New().SetTransport(roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Host == "bbs-api-os.hoyolab.com" {
			if _, ok := r.Context().Deadline(); !ok {
				t.Fatal("optional profile request has no deadline")
			}
			<-r.Context().Done()
			return nil, r.Context().Err()
		}
		if _, ok := r.Context().Deadline(); ok {
			t.Error("profile timeout leaked into daily request")
		}
		return &http.Response{StatusCode: 200, Status: "200 OK", Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"retcode":0,"data":{}}`)), Request: r}, nil
	}))}
	daily := &DailyHoyolab{Referer: "https://act.hoyolab.com/", API: DailyAPI{Endpoint: "https://example.invalid", Info: "/info"}}
	if _, err := daily.GetAccountUserInfo(hoyo); err == nil || !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Fatalf("timeout error=%v", err)
	}
	if _, err := daily.GetCheckInInfo(hoyo); err != nil {
		t.Fatalf("daily after timeout: %v", err)
	}
}
