package user

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserUpdate struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}
