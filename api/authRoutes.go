package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	authUsecase "github.com/semirm-dev/go-common/api/auth"
	"github.com/semirm-dev/go-common/api/domain"
	"github.com/semirm-dev/go-common/api/infra"
	"github.com/semirm-dev/go-common/storage/cache"
)

// authMiddleware will authenticate request
func (auth *Auth) authMiddleware(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, email, token, err := auth.auth(r)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		currentSession, _ := auth.Cache.Get(&cache.Item{Key: email})
		if string(currentSession) != token {
			w.Write([]byte(fmt.Sprint("no session found for token: ", token)))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// register API endpoint will register user into system
func (auth *Auth) register(w http.ResponseWriter, r *http.Request) {
	user := &domain.User{}

	if err := user.DecodeFromReader(r.Body); err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	// TODO: replace with DI framework
	usecase := &authUsecase.UserAccount{
		Repository: &infra.UserAccountRepository{
			SQL: auth.SQL,
		},
		Token: auth.Token,
		Session: &infra.Session{
			Cache: auth.Cache,
		},
	}

	response := usecase.Register(user)

	auth.OnSuccess(response.ToBytes(), w)
}

// login API endpoint will login user into system
func (auth *Auth) login(w http.ResponseWriter, r *http.Request) {
	user := &domain.User{}

	if err := user.DecodeFromReader(r.Body); err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	// TODO: replace with DI framework
	usecase := &authUsecase.UserAccount{
		Repository: &infra.UserAccountRepository{
			SQL: auth.SQL,
		},
		Token: auth.Token,
		Session: &infra.Session{
			Cache: auth.Cache,
		},
	}

	response := usecase.Login(user)

	auth.OnSuccess(response.ToBytes(), w)
}

// logout API endpoint will logout user from system
func (auth *Auth) logout(w http.ResponseWriter, r *http.Request) {
	_, email, _, err := auth.auth(r)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	// TODO: replace with DI framework
	usecase := &authUsecase.UserAccount{
		Token: auth.Token,
		Session: &infra.Session{
			Cache: auth.Cache,
		},
	}

	response := usecase.Logout(email)

	auth.OnSuccess(response.ToBytes(), w)
}

// auth is helper function to extract user id and email from request
func (auth *Auth) auth(r *http.Request) (int, string, string, error) {
	token := r.Header.Get("auth")

	claims, valid := auth.Token.ValidateAndExtract(token)
	if claims == nil || !valid {
		return 0, "", "", errors.New(fmt.Sprint("failed to validate and extract token: ", token))
	}

	idClaim := fmt.Sprint(claims.Fields["id"])

	id, err := strconv.Atoi(idClaim)
	if err != nil {
		return 0, "", "", errors.New(fmt.Sprint("failed to extract user id: ", idClaim))
	}

	email := fmt.Sprint(claims.Fields["email"])

	return id, email, token, nil
}
