package user

import (
	"fmt"

	"github.com/semirm-dev/go-common/crypto"
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

	// TODO: send email with reset token

	return response
}

// UpdatePassword usecase
func (usecase *PasswordResetUsecase) UpdatePassword(passwordReset *PasswordReset) *PasswordUpdateResponse {
	response := &PasswordUpdateResponse{}

	email, err := usecase.Repository.GetByToken(passwordReset.Token)
	if err != nil || email == "" {
		response.ErrorMessage = fmt.Sprint("failed to get email from password reset token: ", err)
		return response
	}

	argon := crypto.NewArgon2()
	argon.Plain = passwordReset.Password

	if err := argon.Hash(); err != nil {
		response.ErrorMessage = fmt.Sprint("failed to hash password: ", err)
		return response
	}

	password := argon.Hashed

	if err := usecase.Repository.UpdatePassword(email, password); err != nil {
		response.ErrorMessage = fmt.Sprintf("update password failed [%s]: %s", email, err)
		return response
	}

	if err := usecase.Repository.DeleteToken(passwordReset.Token); err != nil {
		response.ErrorMessage = fmt.Sprint("failed to delete password reset token: ", err)
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
