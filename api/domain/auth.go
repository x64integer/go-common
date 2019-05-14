package domain

import (
	"io"
	"time"
)

// User entity for authentication
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

// PasswordReset entity
type PasswordReset struct {
	Email string `json:"email"`
}

// DecodeFromReader will decode User entity
func (user *User) DecodeFromReader(body io.Reader) error {
	return decodeFromReader(user, body)
}

// DecodeFromReader will decode PasswordReset entity
func (passwordReset *PasswordReset) DecodeFromReader(body io.Reader) error {
	return decodeFromReader(passwordReset, body)
}
