package films

import (
	"time"
)

// Film модель для документации Swagger
// @description Модель фильма с рейтингом и датой выпуска
type Film struct {
	// @format uuid
	ID string `json:"film_id"`

	// @minLength 1
	// @maxLength 255
	Title string `json:"title" binding:"required"`

	Description string `json:"description"`

	// @minimum 0
	// @maximum 10
	Rating float64 `json:"rating"`

	// @format date
	ReleaseDate time.Time `json:"release_date"`

	// @format date
	CreatedAt time.Time `json:"created_at"`

	// @format date
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateFilm модель для документации Swagger
// @description Модель фильма с необходимым базисом для обновления
type UpdateFilm struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Rating      *float64   `json:"rating"`
	ReleaseDate *time.Time `json:"release_date"`
}

// UserFilm модель для хранения UUID пользователей и фильмов
// @description Модель, в которой хранятся UUID пользователей и фильмов, связанных с ним
type UserFilm struct {
	// Уникальный идентификатор пользователя
	// @format uuid
	UserID string `json:"user_id"`

	// Уникальный идентификатор фильма
	// @format uuid
	FilmID string `json:"film_id"`
}
