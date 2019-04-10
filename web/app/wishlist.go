package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (a *App) wishlistMessage() (*message, error) {
	u, err := url.ParseRequestURI(apiBaseURL + endopointChannelsHistory)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("token", a.oauthAccessToken)
	q.Set("channel", a.wishlistChannel)
	q.Set("count", "200")
	u.RawQuery = q.Encode()
	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	result := &apiResponse{}
	if json.NewDecoder(res.Body).Decode(result); err != nil {
		return nil, err
	}
	attachments := []*attachment{}
	for i := range result.Messages {
		history := result.Messages[len(result.Messages)-1-i]
		if history.MessageID == "" {
			continue
		}
		isDone := false
		reactions := []string{}
		for _, reaction := range history.Reactions {
			if reaction.Name == "done" {
				isDone = true
				break
			}
			reactions = append(reactions, ":"+reaction.Name+":")
		}
		if isDone {
			continue
		}
		ts, err := strconv.ParseFloat(history.TS, 32)
		if err != nil {
			return nil, err
		}
		footerText := fmt.Sprintf(
			"<https://%s.slack.com/archives/%s/p%s|link>",
			a.workspace,
			a.wishlistChannel,
			strings.Replace(history.TS, ".", "", 1),
		)
		attachments = append(attachments, &attachment{
			Text:   strings.Join(reactions, " ") + " " + history.Text,
			Footer: footerText,
			TS:     int(ts),
		})
	}
	message := &message{Attachments: attachments}
	return message, nil
}
