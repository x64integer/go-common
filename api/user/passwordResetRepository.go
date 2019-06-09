package user

// PasswordResetRepository for password reset
type PasswordResetRepository interface {
	// CreateOrUpdate password reset token, return created token
	CreateOrUpdate(string) (string, error)
	// UpdatePassword for given user account email
	UpdatePassword(string, string) error
}
