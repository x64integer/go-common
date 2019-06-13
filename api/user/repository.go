package user

// Repository for user account
type Repository interface {
	// Store new user account
	Store(*Account) error
	// GetByEmail user account
	GetByEmail(string) (*Account, error)
	// Activate user account
	Activate(string) error
	// GetByActivationToken user account
	GetByActivationToken(string) (*Account, error)
}
