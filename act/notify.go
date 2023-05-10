package act

import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

func (hoyo *Hoyolab) NotifyMessage(message string) error {
	if hoyo.Notify.Token == "" || message == "" {
		return nil
	}

	raw, err := resty.New().R().
		SetHeaders(map[string]string{"Authorization": fmt.Sprintf("Bearer %s", hoyo.Notify.Token)}).
		SetFormData(map[string]string{"message": message}).
		Post("https://notify-api.line.me/api/notify")
	if err != nil {
		return fmt.Errorf("notify::%+v", err)
	}
	if raw.StatusCode() != 200 {
		return fmt.Errorf("notify::%s", raw.Status())
	}
	log.Println("NotifyMessage Sending...")
	return nil
}
