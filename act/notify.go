package act

import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

func sendNotification(service, endpoint string, request *resty.Request) error {
	response, err := request.Post(endpoint)
	if err != nil {
		return fmt.Errorf("%s::%+v", service, err)
	}
	if response.StatusCode() != 200 {
		return fmt.Errorf("%s::%s", service, response.Status())
	}
	return nil
}

func (hoyo *Hoyolab) NotifyMessage(message string) error {
	if message == "" {
		return nil
	}
	if hoyo.Notify.LINENotify == "" && hoyo.Notify.Discord == "" {
		return nil
	}

	client := resty.New()

	if hoyo.Notify.LINENotify != "" {
		err := sendNotification("line", "https://notify-api.line.me/api/notify", client.R().
			SetHeader("Authorization", fmt.Sprintf("Bearer %s", hoyo.Notify.LINENotify)).
			SetFormData(map[string]string{"message": message}))
		if err != nil {
			return err
		}
		log.Println("LINENotify Message Sending...")
	}

	if hoyo.Notify.Discord != "" {
		err := sendNotification("Discord", hoyo.Notify.Discord, client.R().SetFormData(map[string]string{"content": message}))
		if err != nil {
			return err
		}
		log.Println("Discord Message Sending...")
	}

	return nil
}
