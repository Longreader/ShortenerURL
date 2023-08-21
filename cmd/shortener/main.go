package main

import (
	"log"
	"net/http"

	"github.com/Longreader/go-shortener-url.git/config"
	"github.com/Longreader/go-shortener-url.git/internal/app"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {

	r := chi.NewRouter()

	config.Setup()

	r.Use(middleware.Recoverer)

	r.Get("/{id:[0-9A-Za-z]+}", app.IDGetHandler)
	r.Post("/", app.ShortenerURLHandler)
	r.Post("/api/shorten", app.APIShortenerURLHandler)

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(config.SERVER_ADDRESS, r))
}
