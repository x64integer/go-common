package user

import "io"

// PasswordReset entity
type PasswordReset struct {
	Email string `json:"email"`
}

// DecodeFromReader will decode password Reset from io.Reader
func (reset *PasswordReset) DecodeFromReader(body io.Reader) error {
	return decodeFromReader(reset, body)
}
