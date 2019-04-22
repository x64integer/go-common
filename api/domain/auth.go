package domain

import (
	"io"
	"time"
)

// User entity for authentication
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"user_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

// DecodeFromReader will decode User entity
func (user *User) DecodeFromReader(body io.Reader) error {
	return decodeFromReader(user, body)
}
