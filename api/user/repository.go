package user

// Repository for user account
type Repository interface {
	// Store new user into storage
	Store(*Account) error
	// GetByEmail user from storage
	GetByEmail(email string) (*Account, error)
}