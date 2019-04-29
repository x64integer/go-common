package auth

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/semirm-dev/go-common/api/domain"
	"github.com/semirm-dev/go-common/api/infra"
	"github.com/semirm-dev/go-common/api/repository"
	"github.com/semirm-dev/go-common/password"

	"github.com/semirm-dev/go-common/jwt"
)

// UserAccount usecase will handle user authentication
type UserAccount struct {
	Repository repository.UserAccount
	*jwt.Token
	*infra.Session
}

// Response for authentication usecase
type Response struct {
	ErrorMessage string `json:"error_message"`
	Email        string `json:"email"`
	Token        string `json:"token"`
}

// Register new user account
func (userAccount *UserAccount) Register(user *domain.User) *Response {
	response := &Response{}

	hashedPassword, err := password.Hash(user.Password)
	if err != nil {
		response.ErrorMessage = fmt.Sprint("failed to hash password: ", err)
		return response
	}

	user.Password = hashedPassword

	if err := userAccount.Repository.Store(user); err != nil {
		response.ErrorMessage = fmt.Sprintf("failed to store new user account [%v]: %s", user, err)
		return response
	}

	token, err := userAccount.loginUser(user)
	if err != nil {
		response.ErrorMessage = fmt.Sprintf("failed to login user [%v]: %s", user, err)
		return response
	}

	response.Email = user.Email
	response.Token = token

	return response
}

// Login user
func (userAccount *UserAccount) Login(user *domain.User) *Response {
	response := &Response{}

	// validate credentials and create login session/token
	existingUser, err := userAccount.Repository.GetByEmail(user.Email)
	if err != nil {
		response.ErrorMessage = fmt.Sprintf("failed to get user by email [%v]: %s", user, err)
		return response
	}

	if existingUser == nil || !password.Valid(existingUser.Password, user.Password) {
		response.ErrorMessage = fmt.Sprint("invalid credentials")
		return response
	}

	token, err := userAccount.loginUser(existingUser)
	if err != nil {
		response.ErrorMessage = fmt.Sprintf("failed to login user [%v]: %s", existingUser, err)
		return response
	}

	response.Token = token

	return response
}

// Logout user
func (userAccount *UserAccount) Logout(email string) *Response {
	response := &Response{}

	if err := userAccount.Session.Destroy(email); err != nil {
		response.ErrorMessage = fmt.Sprintf("failed to destroy user session [%s]: %s", email, err)
	}

	response.Email = email

	return response
}

// ToBytes will marshal Response to []byte
func (response *Response) ToBytes() []byte {
	b, err := json.Marshal(response)
	if err != nil {
		return []byte(fmt.Sprintf("failed to marshal Response: %v, err: %s", response, err.Error()))
	}

	return b
}

// loginUser is helper function to create user token and session
func (userAccount *UserAccount) loginUser(user *domain.User) (string, error) {
	if err := userAccount.Token.Generate(&jwt.Claims{
		Expiration: time.Hour * 24,
		Fields: map[string]interface{}{
			"email": user.Email,
		},
	}); err != nil {
		return "", err
	}

	token := userAccount.Token.Content

	if err := userAccount.Session.Create(user.Email, token); err != nil {
		return "", err
	}

	return token, nil
}
