package app

import (
	"github.com/Longreader/go-shortener-url.git/internal/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Handler struct {
	Store   *storage.Storage
	BaseURL string
}

func NewHandler(store *storage.Storage, url string) *Handler {
	return &Handler{
		Store:   store,
		BaseURL: url,
	}
}
func (h *Handler) InitRouter() chi.Router {

	r := chi.NewRouter()

	r.Use(middleware.Compress(5))
	r.Use(DecompresMiddleware)
	r.Use(middleware.Recoverer)

	r.Get("/{id:[0-9A-Za-z]+}", h.IDGetHandler)
	r.Post("/", h.ShortenerURLHandler)
	r.Post("/api/shorten", h.APIShortenerURLHandler)

	return r
}
