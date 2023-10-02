package middlewares

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Longreader/go-shortener-url.git/internal/app/auth"
	"github.com/Longreader/go-shortener-url.git/internal/tools"
	"github.com/google/uuid"
)

func (m *Middlewares) UserCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("USER")
		if errors.Is(err, http.ErrNoCookie) || len(cookie.Value) < 16 {
			m.setUser(next, w, r, m.createNewUser(w))
			return
		}

		user, err := m.a.Load(cookie.Value)
		if errors.Is(err, auth.ErrUnauthorized) {
			user = m.createNewUser(w)
		} else if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		m.setUser(next, w, r, user)
	})
}

func (m *Middlewares) setUser(next http.Handler, w http.ResponseWriter, r *http.Request, user uuid.UUID) {
	ctx := context.WithValue(r.Context(), tools.ContextKey("userID"), user)
	next.ServeHTTP(w, r.WithContext(ctx))
}

func (m *Middlewares) createNewUser(w http.ResponseWriter) uuid.UUID {
	user, signed := m.a.Gen()

	cookie := http.Cookie{
		Name:     "USER",
		Value:    signed,
		HttpOnly: true,
		Expires:  time.Now().Add(365 * 24 * time.Hour),
		Path:     "/",
	}

	http.SetCookie(w, &cookie)

	return user
}
