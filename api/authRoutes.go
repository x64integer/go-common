package api

import (
	"fmt"
	"net/http"
	"strings"

	authUsecase "github.com/x64integer/go-common/api/auth"
	"github.com/x64integer/go-common/api/domain"
	"github.com/x64integer/go-common/api/infra"
	"github.com/x64integer/go-common/storage/cache"
)

// authMiddleware will authenticate request
func (auth *Auth) authMiddleware(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("auth")

		if strings.TrimSpace(token) == "" {
			w.Write([]byte("invalid token"))
			return
		}

		claims, valid := auth.Token.ValidateAndExtract(token)
		if claims == nil || !valid {
			w.Write([]byte(fmt.Sprint("failed to validate and extract token: ", token)))
			return
		}

		email := fmt.Sprint(claims.Fields["email"])

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
	token := r.Header.Get("auth")

	claims, valid := auth.Token.ValidateAndExtract(token)
	if claims == nil || !valid {
		w.Write([]byte(fmt.Sprint("failed to validate and extract token: ", token)))
		return
	}

	email := fmt.Sprint(claims.Fields["email"])

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
