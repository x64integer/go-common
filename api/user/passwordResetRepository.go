package user

// PasswordResetRepository for password reset
type PasswordResetRepository interface {
	// CreateOrUpdate password reset token, return created token
	CreateOrUpdate(string) (string, error)
	// UpdatePassword for given user account email
	UpdatePassword(string, string) error
	// GetByToken will get email for a given password reset token
	GetByToken(string) (string, error)
	// DeleteToken will delete password reset token
	DeleteToken(string) error
}
