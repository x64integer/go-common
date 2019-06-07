package user

import (
	"io"
	"time"
)

// Account entity for user authentication
type Account struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

// DecodeFromReader will decode user account from io.Reader
func (account *Account) DecodeFromReader(body io.Reader) error {
	return decodeFromReader(account, body)
}
