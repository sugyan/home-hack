package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/sugyan/home-hack/functions/weather"
)

// Response type
type Response struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

func main() {
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
		Text:         result.Description.Text,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
