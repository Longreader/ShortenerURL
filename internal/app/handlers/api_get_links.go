package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Longreader/go-shortener-url.git/internal/app/auth"
	"github.com/Longreader/go-shortener-url.git/internal/repository"
)

type APIGetUserURLsHandlerResponse struct {
	ShortURL    repository.URL `json:"short_url"`
	OriginalURL repository.URL `json:"original_url"`
}

func (h *Handler) APIGetUserURLsHandler(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetUser(r.Context())
	if err != nil {
		log.Printf("unable to parse user uuid: %v", err)
		h.httpJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	links, err := h.st.GetAll(r.Context(), user)
	if err != nil {
		h.httpJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	response := make([]APIGetUserURLsHandlerResponse, 0)
	for _, link := range links {
		response = append(response, APIGetUserURLsHandlerResponse{
			ShortURL:    h.genShortLink(link.ID),
			OriginalURL: link.URL,
		})
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("unable to marshal response: %v", err)
		h.httpJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	if len(response) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		_, err = w.Write(responseJSON)
		if err != nil {
			log.Printf("write failed: %v", err)
		}
	}
}
