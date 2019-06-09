package user

import (
	"io"
	"time"
)

const (
	// Activated user account status
	Activated = 1
	// Deactivated user account status
	Deactivated = 0
)

// Account entity for user authentication
type Account struct {
	ID              int       `json:"id"`
	Username        string    `json:"username"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	Status          int       `json:"status"`
	ActivationToken string    `json:"activation_token"`
	CreatedAt       time.Time `json:"created_at"`
}

// DecodeFromReader will decode user account from io.Reader
func (account *Account) DecodeFromReader(body io.Reader) error {
	return decodeFromReader(account, body)
}
