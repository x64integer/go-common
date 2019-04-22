package repository

import "github.com/x64integer/go-common/api/domain"

// UserAccount repository
type UserAccount interface {
	// Store new user into storage
	Store(*domain.User) error
	// GetByEmail user from storage
	GetByEmail(email string) (*domain.User, error)
}
