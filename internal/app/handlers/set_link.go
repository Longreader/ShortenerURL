package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/Longreader/go-shortener-url.git/internal/app/auth"
	"github.com/Longreader/go-shortener-url.git/internal/repository"
)

func (h *Handler) ShortenerURLHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	b, err := io.ReadAll(r.Body)
	if err != nil || len(b) == 0 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user, err := auth.GetUser(r.Context())
	if err != nil {
		h.httpJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	id, err := h.st.Set(r.Context(), string(b), user)
	if errors.Is(err, repository.ErrURLAlreadyExists) {
		w.WriteHeader(http.StatusConflict)
	} else if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	url := h.genShortLink(id)

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(url))
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

}
