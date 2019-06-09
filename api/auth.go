package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/semirm-dev/go-common/api/user"
	"github.com/semirm-dev/go-common/mail"
	"github.com/semirm-dev/go-common/storage/cache"

	"github.com/semirm-dev/go-common/jwt"
)

// Authenticatable contract
// Such empty definition is used to satisfy reflection dependencies only
// It might change in the future, if such need arises
type Authenticatable interface{}

// Auth configuration
type Auth struct {
	// required
	*jwt.Token
	CacheClient cache.Service

	// confirm email on registration
	RequireConfirmation bool

	// optional, for /register, /login, /logout routes
	UserAccountRepository user.Repository
	// optional, for /password/reset, /password/reset/{token}, /password/update routes
	PasswordResetRepository user.PasswordResetRepository

	// optional
	RegisterPath            string
	LoginPath               string
	LogoutPath              string
	ConfirmRegistrationPath string

	// optional
	PasswordResetRequestPath string
	PasswordResetFormPath    string
	PasswordResetPath        string
	// callback to run from password reset request (click on password reset generated link)
	PasswordResetCallback func(http.ResponseWriter, *http.Request)

	// optional
	MiddlewareFunc func(http.HandlerFunc) http.Handler
	OnError        func(error, http.ResponseWriter)
	OnSuccess      func([]byte, http.ResponseWriter)

	// TODO: Implement customizable Auth entity
	Entity Authenticatable

	mailer *mail.Client
}

// entityField is helper struct to hold information/data from extracted auth Entity (Authenticatable)
type entityField struct {
	authKey   string
	authValue interface{}
	authType  interface{}
	authTable string
}

// apply will setup auth (register, login, logout routes, validate requirements for Auth)
func (auth *Auth) apply(handler Handler) {
	if auth.Token == nil || auth.CacheClient == nil {
		logrus.Fatal("either auth.Token or auth.CacheClient (or both) is not provided")
	}

	auth.defaults()

	if auth.UserAccountRepository != nil {
		handler.HandleFunc(auth.RegisterPath, func(w http.ResponseWriter, r *http.Request) {
			auth.register(w, r)
		}, "POST")

		handler.HandleFunc(auth.LoginPath, func(w http.ResponseWriter, r *http.Request) {
			auth.login(w, r)
		}, "POST")

		handler.Handle(auth.LogoutPath, auth.MiddlewareFunc(func(w http.ResponseWriter, r *http.Request) {
			auth.logout(w, r)
		}), "GET")

		if auth.RequireConfirmation {
			auth.mailer = &mail.Client{
				Sender: mail.DefaultSMTP(),
			}

			handler.HandleFunc(auth.ConfirmRegistrationPath, func(w http.ResponseWriter, r *http.Request) {
				auth.confirmRegistration(w, r)
			}, "GET")
		}

		logrus.Infof(
			"registered auth routes: register -> %v, login -> %v, logout -> %v, account confirmation -> %v",
			auth.RegisterPath,
			auth.LoginPath,
			auth.LogoutPath,
			auth.ConfirmRegistrationPath,
		)
	}

	if auth.PasswordResetRepository != nil {
		handler.HandleFunc(auth.PasswordResetRequestPath, func(w http.ResponseWriter, r *http.Request) {
			auth.createResetToken(w, r)
		}, "POST")

		handler.HandleFunc(auth.PasswordResetFormPath, func(w http.ResponseWriter, r *http.Request) {
			auth.passwordResetForm(w, r)
		}, "GET")

		handler.HandleFunc(auth.PasswordResetPath, func(w http.ResponseWriter, r *http.Request) {
			auth.updatePassword(w, r)
		}, "POST")

		logrus.Infof(
			"registered password reset routes: token request -> %v, reset form -> %v, update password -> %v",
			auth.PasswordResetRequestPath,
			auth.PasswordResetFormPath,
			auth.PasswordResetPath,
		)
	}
}

// Middleware will authenticate request
func (auth *Auth) Middleware(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, email, token, err := auth.Extract(r)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		currentSession, _ := auth.CacheClient.Get(&cache.Item{Key: email})
		if string(currentSession) != token {
			w.Write([]byte(fmt.Sprint("no session found for token: ", token)))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Extract authentication claims from request
func (auth *Auth) Extract(r *http.Request) (int, string, string, error) {
	token := r.Header.Get("auth")

	if strings.TrimSpace(token) == "" {
		return 0, "", "", errors.New(fmt.Sprint("missing token field"))
	}

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

// register API endpoint will register user into system
func (auth *Auth) register(w http.ResponseWriter, r *http.Request) {
	account := &user.Account{}

	if err := account.DecodeFromReader(r.Body); err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	authUsecase := &user.AuthUsecase{
		RequireConfirmation: auth.RequireConfirmation,
		Repository:          auth.UserAccountRepository,
		Token:               auth.Token,
		Session: &user.Session{
			Cache: auth.CacheClient,
		},
		Mailer: auth.mailer,
	}

	response := authUsecase.Register(account)

	auth.OnSuccess(response.ToBytes(), w)
}

// login API endpoint will login user into system
func (auth *Auth) login(w http.ResponseWriter, r *http.Request) {
	account := &user.Account{}

	if err := account.DecodeFromReader(r.Body); err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	authUsecase := &user.AuthUsecase{
		Repository: auth.UserAccountRepository,
		Token:      auth.Token,
		Session: &user.Session{
			Cache: auth.CacheClient,
		},
	}

	response := authUsecase.Login(account)

	auth.OnSuccess(response.ToBytes(), w)
}

// logout API endpoint will logout user from system
func (auth *Auth) logout(w http.ResponseWriter, r *http.Request) {
	_, email, _, err := auth.Extract(r)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	authUsecase := &user.AuthUsecase{
		Token: auth.Token,
		Session: &user.Session{
			Cache: auth.CacheClient,
		},
	}

	response := authUsecase.Logout(email)

	auth.OnSuccess(response.ToBytes(), w)
}

// confirmRegistration API endpoint will confirm user account registration
func (auth *Auth) confirmRegistration(w http.ResponseWriter, r *http.Request) {
	account := &user.Account{}

	if err := account.DecodeFromReader(r.Body); err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	authUsecase := &user.AuthUsecase{
		Repository: auth.UserAccountRepository,
		Token:      auth.Token,
		Session: &user.Session{
			Cache: auth.CacheClient,
		},
	}

	response := authUsecase.ConfirmRegistration(account)

	auth.OnSuccess(response.ToBytes(), w)
}

// createResetToken API endpoint will create password reset token
func (auth *Auth) createResetToken(w http.ResponseWriter, r *http.Request) {
	passwordReset := &user.PasswordReset{}

	if err := passwordReset.DecodeFromReader(r.Body); err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	passwordResetUsecase := &user.PasswordResetUsecase{
		Repository: auth.PasswordResetRepository,
	}

	response := passwordResetUsecase.CreateResetToken(passwordReset.Email)

	w.Write(response.ToBytes())
}

// passwordResetForm API endpoint will show password reset form
func (auth *Auth) passwordResetForm(w http.ResponseWriter, r *http.Request) {
	auth.PasswordResetCallback(w, r)
}

// updatePassword API endpoint will update password
func (auth *Auth) updatePassword(w http.ResponseWriter, r *http.Request) {
	passwordReset := &user.PasswordReset{}

	if err := passwordReset.DecodeFromReader(r.Body); err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	passwordResetUsecase := &user.PasswordResetUsecase{
		Repository: auth.PasswordResetRepository,
	}

	response := passwordResetUsecase.UpdatePassword(passwordReset)

	w.Write(response.ToBytes())
}

// onError default callback
func onError(err error, w http.ResponseWriter) {
	log.Println(err)
}

// onSuccess default callback
func onSuccess(payload []byte, w http.ResponseWriter) {
	w.Write(payload)
}

// defaults is helper function to apply default values
func (auth *Auth) defaults() {
	if auth.OnError == nil {
		auth.OnError = onError
	}

	if auth.OnSuccess == nil {
		auth.OnSuccess = onSuccess
	}

	if auth.MiddlewareFunc == nil {
		auth.MiddlewareFunc = auth.Middleware
	}

	if strings.TrimSpace(auth.RegisterPath) == "" {
		auth.RegisterPath = "/register"
	}

	if strings.TrimSpace(auth.LoginPath) == "" {
		auth.LoginPath = "/login"
	}

	if strings.TrimSpace(auth.LogoutPath) == "" {
		auth.LogoutPath = "/logout"
	}

	if strings.TrimSpace(auth.ConfirmRegistrationPath) == "" {
		auth.ConfirmRegistrationPath = "/account/confirm/{token}"
	}

	if strings.TrimSpace(auth.PasswordResetRequestPath) == "" {
		auth.PasswordResetRequestPath = "/password/reset"
	}

	if strings.TrimSpace(auth.PasswordResetFormPath) == "" {
		auth.PasswordResetFormPath = "/password/reset/{token}"
	}

	if strings.TrimSpace(auth.PasswordResetPath) == "" {
		auth.PasswordResetPath = "/password/update"
	}
}

// extractEntity is helper function to extract auth entity fields and tags
func (auth *Auth) extractEntity(entityToExtract interface{}) []*entityField {
	var entityValue reflect.Value
	var fields []*entityField
	entityKind := reflect.ValueOf(entityToExtract).Kind()

	if entityKind == reflect.Ptr {
		entityValue = reflect.ValueOf(entityToExtract).Elem()
	} else {
		entityValue = reflect.ValueOf(entityToExtract)
	}

	for i := 0; i < entityValue.NumField(); i++ {
		vField := entityValue.Field(i)
		tField := entityValue.Type().Field(i)

		fieldKey := tField.Tag.Get("auth")
		fieldType := tField.Tag.Get("auth_type")
		var fieldValue interface{}

		switch vField.Kind() {
		case reflect.String:
			fieldValue = fmt.Sprint(vField.String())
		case reflect.Int:
			fieldValue = vField.Int()
		case reflect.Float32, reflect.Float64:
			fieldValue = vField.Float()
		}

		fields = append(fields, &entityField{
			authKey:   fieldKey,
			authValue: fieldValue,
			authType:  fieldType,
			authTable: entityValue.Type().Name(),
		})
	}

	return fields
}

// handleFunc is helper function to setup route, map request payload to auth entity
func (auth *Auth) handleFunc(
	w http.ResponseWriter,
	r *http.Request,
	callback func(entity []*entityField) ([]byte, error),
) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		auth.OnError(err, w)
		return
	}

	if err := json.Unmarshal(b, auth.Entity); err != nil {
		auth.OnError(err, w)
		return
	}

	fields := auth.extractEntity(auth.Entity)

	resp, err := callback(fields)
	if err != nil {
		auth.OnError(err, w)
		return
	}

	auth.OnSuccess(resp, w)
}
