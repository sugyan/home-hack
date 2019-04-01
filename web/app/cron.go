package app

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

const endopointChannelsHistory = "/channels.history"

func (a *App) cronWeatherHandler(w http.ResponseWriter, r *http.Request) *appError {
	message, err := a.weatherMessage()
	if err != nil {
		return &appError{err, "failed to fetch weather"}
	}
	message.Channel = a.weather["CHANNEL"]
	message.UserName = a.weather["USERNAME"]
	message.IconEmoji = a.weather["ICONEMOJI"]

	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(message); err != nil {
		return &appError{err, "failed to encode to json"}
	}
	resp, err := http.Post(a.webhookURL.String(), "application/json", buf)
	if err := json.NewEncoder(buf).Encode(message); err != nil {
		return &appError{err, "failed to post message"}
	}
	defer resp.Body.Close()
	return nil
}

func (a *App) cronReminderHandler(w http.ResponseWriter, r *http.Request) *appError {
	u, err := url.ParseRequestURI(apiBaseURL + endopointChannelsHistory)
	if err != nil {
		return &appError{err, "failed to parse URL"}
	}
	q := url.Values{}
	q.Set("token", a.oauthAccessToken)
	q.Set("channel", a.remindChannel)
	q.Set("count", "200")
	u.RawQuery = q.Encode()
	res, err := http.Get(u.String())
	if err != nil {
		return &appError{err, "failed to get histories"}
	}
	defer res.Body.Close()

	result := &apiResponse{}
	if json.NewDecoder(res.Body).Decode(result); err != nil {
		return &appError{err, "failed to decode json"}
	}
	for _, history := range result.Messages {
		if history.MessageID != "" && len(history.Reactions) == 0 {
			log.Printf("%v", history.Text)
		}
	}
	return nil
}
