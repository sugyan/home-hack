package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	apiBaseURL = "https://slack.com/api"

	endopointChannelsHistory = "/channels.history"
)

type message struct {
	ResponseType string        `json:"response_type,omitempty"`
	Channel      string        `json:"channel,omitempty"`
	UserName     string        `json:"username,omitempty"`
	IconEmoji    string        `json:"icon_emoji,omitempty"`
	Text         string        `json:"text"`
	Attachments  []*attachment `json:"attachments"`
}

type attachment struct {
	AuthorIcon string `json:"author_icon,omitempt"`
	AuthorName string `json:"author_name,omitempt"`
	Text       string `json:"text,omitempty"`
	TS         int    `json:"ts,omitempty"`
	Footer     string `json:"footer,omitempty"`
}

type apiResponse struct {
	OK       bool              `json:"ok"`
	Messages []*historyMessage `json:"messages"`
	HasMore  bool              `json:"has_more"`
	Error    string            `json:"error"`
}

type historyMessage struct {
	MessageID string      `json:"client_msg_id"`
	Text      string      `json:"text"`
	TS        string      `json:"ts"`
	Reactions []*reaction `json:"reactions"`
}

type reaction struct {
	Name string `json:"name"`
}

func (a *App) sendMessage(message *message) error {
	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(message); err != nil {
		return err
	}
	resp, err := http.Post(a.webhookURL.String(), "application/json", buf)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf(resp.Status)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Printf("sent message: %v", string(b))

	return nil
}
