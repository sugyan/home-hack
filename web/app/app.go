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
	r                http.Handler
	oauthAccessToken string
	weather          map[string]string
	wishlistChannel  string
	workspace        string
	webhookURL       *url.URL
}

type appError struct {
	err     error
	message string
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		log.Printf("failed to process request: %s [%s]", err.message, err.err.Error())
	}
}

// NewApp function
func NewApp(environ []string) (*App, error) {
	app := &App{
		weather: map[string]string{},
	}
	for _, env := range environ {
		kv := strings.SplitN(env, "=", 2)
		switch kv[0] {
		case "OAUTH_ACCESS_TOKEN":
			app.oauthAccessToken = kv[1]
		case "WEATHER_CITY", "WEATHER_CHANNEL", "WEATHER_USERNAME", "WEATHER_ICONEMOJI":
			app.weather[strings.TrimPrefix(kv[0], "WEATHER_")] = kv[1]
		case "WISHLIST_CHANNEL":
			app.wishlistChannel = kv[1]
		case "WORKSPACE":
			app.workspace = kv[1]
		case "WEBHOOK_URL":
			u, err := url.ParseRequestURI(kv[1])
			if err != nil {
				return nil, err
			}
			app.webhookURL = u
		}
	}
	r := mux.NewRouter()
	// TODO: Verifying requests from Slack
	// https://api.slack.com/docs/verifying-requests-from-slack
	r.Handle("/slash/weather", appHandler(app.slashWeatherHandler))
	r.Handle("/slash/wishlist", appHandler(app.slashWishlistHandler))
	// TODO: Verifying requests from GAE
	r.Handle("/cron/weather", appHandler(app.cronWeatherHandler))
	r.Handle("/cron/wishlist", appHandler(app.cronWishlistHandler))

	app.r = r
	return app, nil
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.r.ServeHTTP(w, r)
}
