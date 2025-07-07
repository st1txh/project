package user

import (
	"github.com/gofrs/uuid"
	"time"
)

// User модель для документации Swagger
// @description Модель пользователя со всеми его данными
type User struct {

	// @format uuid
	ID string `json:"id"`

	// @minLength 1
	// @maxLength 255
	Name string `json:"name" binding:"required"`

	// @minLength 1
	// @maxLength 255
	Email string `json:"email" binding:"required,email"`

	// @format date
	DateOfBirth time.Time `json:"date_of_birth" binding:"required"`
	Gender      string    `json:"gender" binding:"required"`

	// @format date
	CreatedAt time.Time `json:"created_at"`

	// @format date
	UpdatedAt time.Time `json:"updated_at"`

	// @format uuid
	FilmUUID []uuid.UUID `json:"film_id"`
}

// Update модель для аутентификации пользователя
// @description Модель пользователя с данными, необходимыми для обновления
type Update struct {
	// Полное ФИО пользователя
	// @Example "Иванов Иван Иванович"
	// @MinLength 2
	// @MaxLength 100
	Name *string `json:"name"`

	// Электронная почта пользователя
	// @Example "testemail@example.com"
	// @Format email
	Email *string `json:"email"`

	// Информация о дате рождения пользователя
	// @Example "2000.01.01"
	// @Format date
	DateOfBirth *time.Time `json:"date_of_birth"`

	// Пол пользователя
	// @Enum "М" "Ж"
	// @Format string
	// @MaxLength 1
	Gender *string `json:"gender"`

	// Уникальный идентификатор фильма, с которым связан пользователь
	// @Example "1111a111-2b2b-3333-444d-55555555eee5"
	// @DFormat uuid
	FilmUUID []uuid.UUID `json:"film_id"`
}
