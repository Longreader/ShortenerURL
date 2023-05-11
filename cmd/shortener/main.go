package main

import (
	"log"
	"net/http"

	"github.com/Longreader/go-shortener-url.git/internal/app"
	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/{id:[0-9A-Za-z]+}", app.IDGetHandler)
	r.HandleFunc("/", app.ShortenerURLHandler)

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
