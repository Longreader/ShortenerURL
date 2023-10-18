package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Longreader/go-shortener-url.git/internal/app/auth"
	"github.com/Longreader/go-shortener-url.git/internal/repository"
	"github.com/sirupsen/logrus"
)

type (
	APISetShortenerURLsHandlerRequest struct {
		CorrelationID string         `json:"correlation_id"`
		OriginalURL   repository.URL `json:"original_url"`
	}

	APISetShortenerURLsHandlerResponse struct {
		CorrelationID string         `json:"correlation_id"`
		ShortURL      repository.URL `json:"short_url"`
	}
)

func (h *Handler) APISetShortenerURLsHandler(w http.ResponseWriter, r *http.Request) {

	b, err := io.ReadAll(r.Body)
	if err != nil || len(b) == 0 {
		h.httpJSONError(w, "Bad request", http.StatusBadRequest)
	}

	request := make([]APISetShortenerURLsHandlerRequest, 0)

	err = json.Unmarshal(b, &request)
	if err != nil {
		h.httpJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}

	user, err := auth.GetUser(r.Context())
	if err != nil {
		logrus.Printf("unable to parse user uuid: %v", err)
		h.httpJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	response := make([]APISetShortenerURLsHandlerResponse, 0)

	for _, link := range request {
		id, err := h.st.Set(r.Context(), link.OriginalURL, user)
		if !errors.Is(err, repository.ErrURLAlreadyExists) && err != nil {
			h.httpJSONError(w, "Server error", http.StatusInternalServerError)
			return
		}
		response = append(response, APISetShortenerURLsHandlerResponse{
			CorrelationID: link.CorrelationID,
			ShortURL:      h.genShortLink(id),
		})
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		h.httpJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(responseJSON)
	if err != nil {
		logrus.Printf("write failed: %v", err)
	}
}
