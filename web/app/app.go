package app

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
)

// App type
type App struct {
	r          http.Handler
	weather    map[string]string
	webhookURL *url.URL
}

// NewApp function
func NewApp(environ []string) (*App, error) {
	app := &App{
		weather: map[string]string{},
	}
	for _, env := range environ {
		log.Printf("env: %s", env)
		kv := strings.SplitN(env, "=", 2)
		switch kv[0] {
		case "WEATHER_CITY", "WEATHER_CHANNEL", "WEATHER_USERNAME", "WEATHER_ICONEMOJI":
			app.weather[strings.TrimPrefix(kv[0], "WEATHER_")] = kv[1]
		case "WEBHOOK_URL":
			u, err := url.ParseRequestURI(kv[1])
			if err != nil {
				return nil, err
			}
			app.webhookURL = u
		}
	}
	log.Printf("%v", app.weather)
	r := mux.NewRouter()
	// TODO: Verifying requests from Slack
	// https://api.slack.com/docs/verifying-requests-from-slack
	r.HandleFunc("/slash/weather", app.slashWeatherHandler)
	// TODO: Verifying requests from GAE
	r.HandleFunc("/cron/weather", app.cronWeatherHandler)

	app.r = r
	return app, nil
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.r.ServeHTTP(w, r)
}
