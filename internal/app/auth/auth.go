package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"

	"github.com/Longreader/go-shortener-url.git/config"
	"github.com/google/uuid"
)

var ErrUnauthorized = errors.New("user unauthorized")

type Auth struct {
	cfg config.Config
}

func NewAuth(cfg config.Config) Auth {
	return Auth{
		cfg: cfg,
	}
}

func (a Auth) Load(s string) (user uuid.UUID, err error) {
	payload, err := base64.StdEncoding.DecodeString(s)
	if err != nil || len(payload) < 16 {
		return uuid.Nil, ErrUnauthorized
	}

	h := hmac.New(sha256.New, a.cfg.CookieKey)
	h.Write(payload[:16])
	sign := h.Sum(nil)

	if !hmac.Equal(sign, payload[16:]) {
		return uuid.Nil, ErrUnauthorized
	}

	user, err = uuid.FromBytes(payload[:16])
	if err != nil {
		return uuid.Nil, ErrUnauthorized
	}

	return user, nil
}

func (a Auth) Gen() (user uuid.UUID, signed string) {
	user = uuid.New()

	b, _ := user.MarshalBinary()

	h := hmac.New(sha256.New, a.cfg.CookieKey)
	h.Write(b)
	sign := h.Sum(nil)

	signed = base64.StdEncoding.EncodeToString(append(b, sign...))

	return user, signed
}
