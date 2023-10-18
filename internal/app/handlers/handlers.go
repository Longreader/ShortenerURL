package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Longreader/go-shortener-url.git/internal/storage"
)

type Handler struct {
	st     storage.Storager
	domain string
}

func NewHandler(st storage.Storager, domanin string) *Handler {
	return &Handler{
		st:     st,
		domain: domanin,
	}
}

func (h *Handler) httpJSONError(w http.ResponseWriter, msg string, code int) {
	jsonError, _ := json.Marshal(
		struct {
			Error string `json:"error"`
		}{
			Error: msg,
		},
	)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	_, err := w.Write(jsonError)
	if err != nil {
		log.Printf("write failed: %v", err)
	}
}

func (h *Handler) genShortLink(id string) string {
	return fmt.Sprintf("%s/%s", h.domain, id)
}
