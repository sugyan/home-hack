package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sugyan/home-hack/web/app"
)

func main() {
	app, err := app.NewApp(os.Environ())
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatal(err)
	}
}
