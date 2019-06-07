package user

// PasswordResetRepository for password reset
type PasswordResetRepository interface {
	// CreateOrUpdate password reset token, return created token
	CreateOrUpdate(string) (string, error)
}
