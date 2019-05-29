package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type HangoutsChatNotifier struct {
	WebhookUrl string `json:"webhook_url"`
}

func (n *HangoutsChatNotifier) Send(message string) error {
	url := n.WebhookUrl

	var jsonStr = []byte(fmt.Sprintf(`{"text":"%s"}`, escapeString(message)))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json.RawMessage(jsonStr)))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil

}
