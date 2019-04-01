package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sugyan/home-hack/functions/weather"
)

func (a *App) slashWeatherHandler(w http.ResponseWriter, r *http.Request) *appError {
	message, err := a.weatherMessage()
	if err != nil {
		return &appError{err, "failed to fetch weather"}
	}
	message.ResponseType = "in_channel"
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(message); err != nil {
		return &appError{err, "failed to encode json"}
	}
	return nil
}

func (a *App) weatherMessage() (*message, error) {
	cityIDstr := a.weather["CITY"]
	cityID, err := strconv.Atoi(cityIDstr)
	if err != nil {
		return nil, err
	}
	result, err := weather.FetchForecast(cityID)
	if err != nil {
		return nil, err
	}
	message := &message{
		Text:        fmt.Sprintf("%s (%s 発表) %s", result.Title, result.PublicTime.Time.Format(time.Kitchen), result.Link),
		Attachments: []*attachment{},
	}
	for _, forecast := range result.Forecasts {
		attachment := &attachment{
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
