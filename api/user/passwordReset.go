package user

import "io"

// PasswordReset entity
type PasswordReset struct {
	Email    string `json:"email"`
	Token    string `json:"token,omitempty"`
	Password string `json:"password,omitempty"`
}

// DecodeFromReader will decode password Reset from io.Reader
func (reset *PasswordReset) DecodeFromReader(body io.Reader) error {
	return decodeFromReader(reset, body)
}
