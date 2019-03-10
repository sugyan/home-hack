package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sugyan/home-hack/functions/weather"
)

// Response type
type Response struct {
	ResponseType string        `json:"response_type"`
	Text         string        `json:"text"`
	Attachments  []*Attachment `json:"attachments"`
}

// Attachment type
type Attachment struct {
	AuthorIcon string `json:"author_icon"`
	AuthorName string `json:"author_name"`
	Text       string `json:"text"`
}

func main() {
	// TODO: Verifying requests from Slack
	// https://api.slack.com/docs/verifying-requests-from-slack
	http.HandleFunc("/slash/weather", weatherHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatal(err)
	}
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	cityIDstr := os.Getenv("WEATHER_CITY")
	cityID, err := strconv.Atoi(cityIDstr)
	if err != nil {
		log.Printf(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	result, err := weather.FetchForecast(cityID)
	if err != nil {
		log.Printf(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	res := &Response{
		ResponseType: "in_channel",
		Text:         fmt.Sprintf("%s (%s 発表) %s", result.Title, result.PublicTime.Time.Format(time.Kitchen), result.Link),
		Attachments:  []*Attachment{},
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
		res.Attachments = append(res.Attachments, attachment)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
