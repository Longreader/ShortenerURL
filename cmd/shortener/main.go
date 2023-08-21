package main

import (
	"log"
	"net/http"

	"os"

	"github.com/Longreader/go-shortener-url.git/internal/app"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {

	var serverAddress = os.Getenv("SERVER_ADDRESS")

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)

	r.Get("/{id:[0-9A-Za-z]+}", app.IDGetHandler)
	r.Post("/", app.ShortenerURLHandler)
	r.Post("/api/shorten", app.APIShortenerURLHandler)

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(serverAddress, r))
}
