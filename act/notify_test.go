package act

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestSendNotification(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if err := request.ParseForm(); err != nil {
			t.Fatalf("ParseForm returned error: %v", err)
		}
		if request.FormValue("content") != "daily reward" {
			t.Fatalf("unexpected form content: %q", request.FormValue("content"))
		}
		response.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	request := resty.New().R().SetFormData(map[string]string{"content": "daily reward"})
	if err := sendNotification("Discord", server.URL, request); err != nil {
		t.Fatalf("sendNotification returned error: %v", err)
	}
}

func TestSendNotificationStatusError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		response.WriteHeader(http.StatusBadGateway)
	}))
	t.Cleanup(server.Close)

	err := sendNotification("Discord", server.URL, resty.New().R())
	if err == nil || !strings.Contains(err.Error(), "Discord::502 Bad Gateway") {
		t.Fatalf("unexpected status error: %v", err)
	}
}
