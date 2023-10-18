package repository

import (
	"github.com/google/uuid"
)

type (
	ID   = string    // Тип для хранения ID сокращенной ссылки.
	URL  = string    // Тип для хранения исходного URL.
	User = uuid.UUID // Тип для хранения ID пользователя.
)

type LinkData struct {
	ID   ID   // ID сокращенной ссылки.
	URL  URL  // Исходный URL.
	User User // Пользователь, которому принадлежит ссылка.
}
