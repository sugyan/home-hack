package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sugyan/home-hack/functions/weather"
)

func (a *App) slashWeatherHandler(w http.ResponseWriter, r *http.Request) {
	message, err := a.weatherMessage()
	if err != nil {
		log.Printf(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	message.ResponseType = "in_channel"
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(message); err != nil {
		log.Printf(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (a *App) cronWeatherHandler(w http.ResponseWriter, r *http.Request) {
	message, err := a.weatherMessage()
	if err != nil {
		log.Printf(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	message.Channel = a.weather["CHANNEL"]
	message.UserName = a.weather["USERNAME"]
	message.IconEmoji = a.weather["ICONEMOJI"]

	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(message); err != nil {
		log.Printf(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	resp, err := http.Post(a.webhookURL.String(), "application/json", buf)
	if err := json.NewEncoder(buf).Encode(message); err != nil {
		log.Printf(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Write([]byte("OK"))
}

func (a *App) weatherMessage() (*Message, error) {
	cityIDstr := a.weather["CITY"]
	cityID, err := strconv.Atoi(cityIDstr)
	if err != nil {
		return nil, err
	}
	result, err := weather.FetchForecast(cityID)
	if err != nil {
		return nil, err
	}
	message := &Message{
		Text:        fmt.Sprintf("%s (%s 発表) %s", result.Title, result.PublicTime.Time.Format(time.Kitchen), result.Link),
		Attachments: []*Attachment{},
	}
	for _, forecast := range result.Forecasts {
		attachment := &Attachment{
			AuthorIcon: forecast.Image.URL,
			AuthorName: fmt.Sprintf("%s (%s)", forecast.DateLabel, forecast.Date),
			Text:       forecast.Telop,
		}
		temperatures := []string{}
		if forecast.Temperature.Max != nil {
			temperatures = append(temperatures, fmt.Sprintf("最高気温: %s℃", forecast.Temperature.Max.Celsius))
		}
		if forecast.Temperature.Min != nil {
			temperatures = append(temperatures, fmt.Sprintf("最低気温: %s℃", forecast.Temperature.Min.Celsius))
		}
		if len(temperatures) > 0 {
			attachment.Text += fmt.Sprintf(" (%s)", strings.Join(temperatures, ", "))
		}
		message.Attachments = append(message.Attachments, attachment)
	}
	return message, nil
}
