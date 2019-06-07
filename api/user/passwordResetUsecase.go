package user

import (
	"fmt"
)

// PasswordResetUsecase for password reset
type PasswordResetUsecase struct {
	Repository PasswordResetRepository
}

// PasswordResetResponse for password reset
type PasswordResetResponse struct {
	ErrorMessage string `json:"error_message"`
	Token        string `json:"token"`
}

// CreateResetToken for password reset usecase
func (usecase *PasswordResetUsecase) CreateResetToken(email string) *PasswordResetResponse {
	response := &PasswordResetResponse{}

	token, err := usecase.Repository.CreateOrUpdate(email)
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
