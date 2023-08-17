package app

import (
	"io"
	"net/http"

	"github.com/Longreader/go-shortener-url.git/internal/shortener" //indirect
	"github.com/Longreader/go-shortener-url.git/internal/storage"
	"github.com/gorilla/mux"
)

var Store = storage.New()

func IDGetHandler(w http.ResponseWriter, r *http.Request) {
	// Handle 'GET' method
	// income ID-shorl url
	// outcome Redirect by full url
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET request are allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)

	// WORK with storage
	// Move id - data url
	// url, ok := storage[vars["id"]]
	fullURL, ok := Store.Get(vars["id"])

	if !ok {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	// log.Println("IDFefaultHandler func()\n\tOK is", ok, "\n\tItem is", fullURL)

	w.Header().Add("Content-Type", "text/html; charset=UTF-8")
	w.Header().Set("Location", string(fullURL))
	w.WriteHeader(307)
}

func ShortenerURLHandler(w http.ResponseWriter, r *http.Request) {
	// Handle POST method
	// income full URL
	// outcome None
	// set up pair 'full url' - 'id' in storage
	if r.Method != http.MethodPost {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	// Read URL
	fullURL, err := io.ReadAll(r.Body)
	// Handle error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// WORK with full url
	shortURL := shortener.RandStringBytes(7)
	for _, ok := Store.Get(shortURL); ok; {
		shortURL := shortener.RandStringBytes(7)
		_, ok = Store.Get(shortURL)
	}
	baseURL := "http://127.0.0.1:8080/"
	// Write shor url to body
	w.WriteHeader(201)
	w.Write([]byte(baseURL + shortURL))
	// Set shorturl-url to store map
	Store.Set(shortURL, string(fullURL))
	// fmt.Println(reflect.TypeOf(b))
}
