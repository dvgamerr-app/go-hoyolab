package act

import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

func (hoyo *Hoyolab) NotifyMessage(message string) error {
	if message == "" {
		return nil
	}

	if hoyo.Notify.LINENotify != "" {
		raw, err := resty.New().R().
			SetHeaders(map[string]string{"Authorization": fmt.Sprintf("Bearer %s", hoyo.Notify.LINENotify)}).
			SetFormData(map[string]string{"message": message}).
			Post("https://notify-api.line.me/api/notify")
		if err != nil {
			return fmt.Errorf("line::%+v", err)
		}
		if raw.StatusCode() != 200 {
			return fmt.Errorf("line::%s", raw.Status())
		}
		log.Println("LINENotify Message Sending...")
	}

	if hoyo.Notify.Discord != "" {
		raw, err := resty.New().R().
			SetFormData(map[string]string{"content": message}).
			Post(hoyo.Notify.Discord)
		if err != nil {
			return fmt.Errorf("Discord::%+v", err)
		}
		if raw.StatusCode() != 200 {
			return fmt.Errorf("Discord::%s", raw.Status())
		}
		log.Println("Discord Message Sending...")
	}

	return nil
}
