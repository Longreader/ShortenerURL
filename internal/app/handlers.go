package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Longreader/go-shortener-url.git/internal/shortener" //indirect
	"github.com/Longreader/go-shortener-url.git/internal/storage"
	"github.com/go-chi/chi/v5"
)

var Store = storage.New()
var baseURL = "http://127.0.0.1:8080/"

func IDGetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET request are allowed", http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	fullURL, ok := Store.Get(id)

	if !ok {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "text/html; charset=UTF-8")
	w.Header().Set("Location", string(fullURL))
	w.WriteHeader(307)
}

func ShortenerURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	// Читаем из тела запроса
	fullURL, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	// Отлавливаем возможную ошибку
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Создание сокращенного URL
	shortURL := shortener.RandStringBytes(7)
	for _, ok := Store.Get(shortURL); ok; {
		shortURL = shortener.RandStringBytes(7)
		_, ok = Store.Get(shortURL)
	}
	// baseURL := "http://127.0.0.1:8080/"

	w.WriteHeader(201)
	w.Write([]byte(baseURL + shortURL))

	Store.Set(shortURL, string(fullURL))
}

func APIShortenerURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	// Чтение JSON из тела запроса
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Bad agent request", http.StatusNotAcceptable)
		return
	}
	defer r.Body.Close()

	docoder := struct {
		URL string `json:"url"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&docoder); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(docoder.URL)

	// Создание сокращенного URL
	shortURL := shortener.RandStringBytes(7)
	for _, ok := Store.Get(shortURL); ok; {
		shortURL = shortener.RandStringBytes(7)
		_, ok = Store.Get(shortURL)
	}
	// baseURL := "http://127.0.0.1:8080/"
	fmt.Println(docoder.URL)
	Store.Set(shortURL, string(docoder.URL))

	// Запись ответа JSON в тело ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	value := struct {
		Result string `json:"result"`
	}{
		Result: baseURL + shortURL,
	}

	buf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buf)
	encoder.Encode(value)
	w.Write(buf.Bytes())
}
