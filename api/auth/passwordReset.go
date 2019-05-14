package auth

import (
	"fmt"

	"github.com/semirm-dev/go-common/api/repository"
)

// PasswordReset usecase
type PasswordReset struct {
	Repository repository.PasswordReset
}

// PasswordResetResponse model
type PasswordResetResponse struct {
	ErrorMessage string `json:"error_message"`
	Token        string `json:"token"`
}

// CreateResetToken for password reset usecase
func (passwordReset *PasswordReset) CreateResetToken(email string) *PasswordResetResponse {
	response := &PasswordResetResponse{}

	token, err := passwordReset.Repository.CreateOrUpdate(email)
	if err != nil {
		response.ErrorMessage = fmt.Sprintf("failed to create password reset token [%s]: %s", email, err)
		return response
	}

	response.Token = token

	return response
}

// ToBytes will marshal Response to []byte
func (response *PasswordResetResponse) ToBytes() []byte {
	return toBytes(response)
}
