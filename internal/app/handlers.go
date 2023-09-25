package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Longreader/go-shortener-url.git/internal/tools"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

func (h *Handler) IDGetHandler(w http.ResponseWriter, r *http.Request) {

	logrus.Debug("IDGetHandler start")
	defer logrus.Debug("IDGetHandler end")

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET request are allowed", http.StatusMethodNotAllowed)
		return
	}

	// h.Store.GetAll()

	id := chi.URLParam(r, "id")

	logrus.Debugf("ID equal %s", id)

	if id == "" {
		http.Error(w, "Bad request or miss id", http.StatusBadRequest)
		return
	}

	fullURL, ok := h.Store.Get(id)

	if !ok {
		http.Error(w, fmt.Sprintf("Wrong ID %s", fullURL), http.StatusBadRequest)
		return
	}

	w.Header().Add("Location", string(fullURL))
	w.Header().Add("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) ShortenerURLHandler(w http.ResponseWriter, r *http.Request) {

	logrus.Debug("ShortenerURLHandler start")
	defer logrus.Debug("ShortenerURLHandler end")
	logrus.Debugf("Base URL %s", h.BaseURL)

	if r.Method != http.MethodPost {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	fullURL, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shortURL := tools.RandStringBytes(7)
	for _, ok := h.Store.Get(shortURL); ok; {
		shortURL = tools.RandStringBytes(7)
		_, ok = h.Store.Get(shortURL)
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte(h.BaseURL + shortURL))

	h.Store.Set(shortURL, string(fullURL))
}

func (h *Handler) APIShortenerURLHandler(w http.ResponseWriter, r *http.Request) {

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

	shortURL := tools.RandStringBytes(7)
	for _, ok := h.Store.Get(shortURL); ok; {
		shortURL = tools.RandStringBytes(7)
		_, ok = h.Store.Get(shortURL)
	}
	h.Store.Set(shortURL, string(docoder.URL))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	value := struct {
		Result string `json:"result"`
	}{
		Result: h.BaseURL + shortURL,
	}

	buf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buf)
	encoder.Encode(value)
	w.Write(buf.Bytes())
}
