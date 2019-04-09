package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// SlackNotifier defines the spec for integrating with Slack
// +k8s:openapi-gen=true
type SlackNotifier struct {
	WebhookUrl string `json:"webhook_url"`
	Channel    string `json:"channel,omitempty"`
	Username   string `json:"username,omitempty"`
	IconEmoji  string `json:"icon_emoji,omitempty"`
}

func (n *SlackNotifier) Send(message string) error {
	url := n.WebhookUrl
	channel := ``
	if n.Channel != `` {
		channel = fmt.Sprintf(`, "channel":"%s"`, n.Channel)
	}
	username := ``
	if n.Username != `` {
		username = fmt.Sprintf(`, "username":"%s"`, n.Username)
	}
	icon_emoji := ``
	if n.IconEmoji != `` {
		icon_emoji = fmt.Sprintf(`, "icon_emoji":"%s"`, n.IconEmoji)
	}

	var jsonStr = []byte(fmt.Sprintf(`{"text":"%s"%s%s%s}`, escapeString(message), channel, username, icon_emoji))
	fmt.Println("message: ", string(jsonStr))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json.RawMessage(jsonStr)))
	req.Header.Set("Content-Type", "application/json")

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
