package middlewares

import (
	"github.com/Longreader/go-shortener-url.git/config"
	"github.com/Longreader/go-shortener-url.git/internal/app/auth"
)

type Middlewares struct {
	cfg config.Config
	a   auth.Auth
}

func NewMiddlewares(cfg config.Config, a auth.Auth) Middlewares {
	return Middlewares{
		cfg: cfg,
		a:   a,
	}
}
