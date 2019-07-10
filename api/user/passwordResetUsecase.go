package user

import (
	"fmt"

	"github.com/semirm-dev/go-common/crypto"
	"github.com/semirm-dev/go-common/mail"
)

// PasswordResetUsecase for password reset
type PasswordResetUsecase struct {
	Repository            PasswordResetRepository
	Mailer                *mail.Client
	ConfirmResetTokenPath string
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

	if err := usecase.sendTokenResetMail(email, token); err != nil {
		response.ErrorMessage = fmt.Sprintf("failed to send reset token confirmation mail [%s][%s]: %s", email, token, err)
		return response
	}

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

// sendTokenResetMail will cosntruct password reset mail body and send it
//
// TODO: parse subject and body from external template
func (usecase *PasswordResetUsecase) sendTokenResetMail(to string, token string) error {
	subject := "Password reset request"
	body := []byte("Click on the link to reset password: <a href=\"http://" + usecase.ConfirmResetTokenPath + token + "\">Reset</a>")

	content := &mail.Content{
		To:      []string{to},
		Subject: subject,
		Body:    body,
	}

	return usecase.Mailer.Send(content)
}
