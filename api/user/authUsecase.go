package user

import (
	"fmt"
	"strings"
	"time"

	"github.com/semirm-dev/go-common/jwt"
	"github.com/semirm-dev/go-common/mail"
	"github.com/semirm-dev/go-common/password"
)

// AuthUsecase will handle user authentication
type AuthUsecase struct {
	RequireConfirmation bool
	Repository
	*jwt.Token
	*Session
	Mailer *mail.Client
}

// AuthResponse for authentication
type AuthResponse struct {
	ErrorMessage string `json:"error_message"`
	ID           int    `json:"id"`
	Email        string `json:"email"`
	Token        string `json:"token"`
}

// Register new user account
func (usecase *AuthUsecase) Register(user *Account) *AuthResponse {
	response := &AuthResponse{}

	hashedPassword, err := password.Hash(user.Password)
	if err != nil {
		response.ErrorMessage = fmt.Sprint("failed to hash password: ", err)
		return response
	}

	user.Password = hashedPassword

	if !usecase.RequireConfirmation {
		user.Status = Activated
	}

	if strings.TrimSpace(user.Email) == "" {
		response.ErrorMessage = "email is required"
		return response
	}

	if strings.TrimSpace(user.Username) == "" {
		user.Username = user.Email
	}

	if err := usecase.Repository.Store(user); err != nil {
		response.ErrorMessage = fmt.Sprintf("failed to store new user account [%v]: %s", user, err)
		return response
	}

	response.ID = user.ID
	response.Email = user.Email

	if usecase.RequireConfirmation {
		content := &mail.Content{
			To:      []string{user.Email},
			Subject: "Confirm registration",
			Body:    []byte("Please confirm email by clicking on the link: <token>"),
		}

		if err := usecase.Mailer.Send(content); err != nil {
			response.ErrorMessage = fmt.Sprintf("failed to send email confirmation [%v]: %s", user, err)
			return response
		}
	} else {
		token, err := usecase.loginUser(user)
		if err != nil {
			response.ErrorMessage = fmt.Sprintf("failed to login user [%v]: %s", user, err)
			return response
		}

		response.Token = token
	}

	return response
}

// Login user
func (usecase *AuthUsecase) Login(user *Account) *AuthResponse {
	response := &AuthResponse{}

	// validate credentials and create login session/token
	existingUser, err := usecase.Repository.GetByEmail(user.Email)
	if err != nil {
		response.ErrorMessage = fmt.Sprintf("failed to get user by email [%v]: %s", user, err)
		return response
	}

	if existingUser == nil || !password.Valid(existingUser.Password, user.Password) {
		response.ErrorMessage = fmt.Sprint("invalid credentials")
		return response
	}

	if existingUser.Status != Activated {
		response.ErrorMessage = fmt.Sprint("account not activated")
		return response
	}

	token, err := usecase.loginUser(existingUser)
	if err != nil {
		response.ErrorMessage = fmt.Sprintf("failed to login user [%v]: %s", existingUser, err)
		return response
	}

	response.Token = token

	return response
}

// Logout user
func (usecase *AuthUsecase) Logout(email string) *AuthResponse {
	response := &AuthResponse{}

	if err := usecase.Session.Destroy(email); err != nil {
		response.ErrorMessage = fmt.Sprintf("failed to destroy user session [%s]: %s", email, err)
	}

	response.Email = email

	return response
}

// ConfirmRegistration for user account
func (usecase *AuthUsecase) ConfirmRegistration(user *Account) *AuthResponse {
	response := &AuthResponse{}

	if err := usecase.Repository.Activate(user.ActivationToken); err != nil {
		response.ErrorMessage = fmt.Sprintf("failed to activate user [%v]: %s", user, err)
		return response
	}

	token, err := usecase.loginUser(user)
	if err != nil {
		response.ErrorMessage = fmt.Sprintf("failed to login user [%v]: %s", user, err)
		return response
	}

	response.Token = token

	return response
}

// ToBytes will marshal Response to []byte
func (response *AuthResponse) ToBytes() []byte {
	return toBytes(response)
}

// loginUser is helper function to create user token and session
func (usecase *AuthUsecase) loginUser(user *Account) (string, error) {
	if err := usecase.Token.Generate(&jwt.Claims{
		Expiration: time.Hour * 24,
		Fields: map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
		},
	}); err != nil {
		return "", err
	}

	token := usecase.Token.Content

	if err := usecase.Session.Create(user.Email, token); err != nil {
		return "", err
	}

	return token, nil
}
