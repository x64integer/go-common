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

// PasswordUpdateResponse for password reset
type PasswordUpdateResponse struct {
	ErrorMessage string `json:"error_message"`
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

// UpdatePassword usecase
func (usecase *PasswordResetUsecase) UpdatePassword(passwordReset *PasswordReset) *PasswordUpdateResponse {
	response := &PasswordUpdateResponse{}

	// TODO: get email from password_reset based on given token, delete doken, update password

	if err := usecase.Repository.UpdatePassword(passwordReset.Email, passwordReset.Password); err != nil {
		response.ErrorMessage = fmt.Sprintf("update password failed [%s]: %s", passwordReset.Email, err)
		return response
	}

	return response
}

// ToBytes will marshal Response to []byte
func (response *PasswordResetResponse) ToBytes() []byte {
	return toBytes(response)
}

// ToBytes will marshal Response to []byte
func (response *PasswordUpdateResponse) ToBytes() []byte {
	return toBytes(response)
}
