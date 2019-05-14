package repository

// PasswordReset repository
type PasswordReset interface {
	// CreateOrUpdate password reset token, return created token
	CreateOrUpdate(string) (string, error)
}
