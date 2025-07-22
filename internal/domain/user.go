package domain

import (
	"github.com/google/uuid"
	"time"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,min=3,max=30,regexp=^[\\w.]+$"`
	Email    string `json:"email"    validate:"required,email,max=150"`
	Password string `json:"password" validate:"required,min=3,max=130,matches=^(?=.*[A-Za-z])(?=.*\\d)[[:alnum:][:punct:]]+$"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
