package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HangoutsChatNotifier struct {
	WebhookUrl string `json:"webhook_url"`
}

func (n *HangoutsChatNotifier) Send(message string) error {
	url := n.WebhookUrl

	var jsonStr = []byte(fmt.Sprintf(`{"text":"%s"}`, escapeString(message)))
	fmt.Println("message: ", string(jsonStr))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json.RawMessage(jsonStr)))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return nil

}
