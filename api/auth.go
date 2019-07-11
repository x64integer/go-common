package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/semirm-dev/go-common/api/user"
	"github.com/semirm-dev/go-common/storage/cache"

	"github.com/semirm-dev/go-common/jwt"
)

const (
	accountConfirm        = "/account/confirm/"
	passwordResetPath     = "/password/reset/"
	passwordResetTemplate = "../go-common/api/tpl/password_reset.gohtml"
)

// Authenticatable contract
//
// Such empty definition is used to satisfy reflection dependencies only.
// It might change in the future, if such need arises
type Authenticatable interface{}

// Auth configuration
type Auth struct {
	// required
	*jwt.Token
	CacheClient cache.Service

	Router

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
	confirmRegistrationPath string
	ServiceURL              string

	// optional
	PasswordResetRequestPath string
	passwordResetFormPath    string
	PasswordResetPath        string
	// callback to run from password reset request (click on password reset generated link)
	PasswordResetCallback func(http.ResponseWriter, *http.Request)

	// optional
	MiddlewareFunc func(http.HandlerFunc) http.Handler
	OnError        func(error, http.ResponseWriter)
	OnSuccess      func([]byte, http.ResponseWriter)

	// TODO: Implement customizable Auth entity
	Entity Authenticatable
}

// entityField is helper struct to hold information/data from extracted auth Entity (Authenticatable)
type entityField struct {
	authKey   string
	authValue interface{}
	authType  interface{}
	authTable string
}

// Apply will setup auth (register, login, logout routes, validate requirements for Auth)
func (auth *Auth) Apply(router Router) {
	if auth.Token == nil || auth.CacheClient == nil {
		logrus.Fatal("either auth.Token or auth.CacheClient (or both) is not provided")
	}

	auth.Router = router

	auth.defaults()

	if auth.UserAccountRepository != nil {
		router.HandleFunc(auth.RegisterPath, auth.register, "POST")
		router.HandleFunc(auth.LoginPath, auth.login, "POST")
		router.Handle(auth.LogoutPath, auth.MiddlewareFunc(auth.logout), "GET")

		if auth.RequireConfirmation {
			router.HandleFunc(auth.confirmRegistrationPath, auth.confirmRegistration, "GET")
		}

		logrus.Infof(
			"auth routes: register -> %v, login -> %v, logout -> %v, account confirmation -> %v",
			auth.RegisterPath,
			auth.LoginPath,
			auth.LogoutPath,
			auth.confirmRegistrationPath,
		)
	}

	if auth.PasswordResetRepository != nil {
		router.HandleFunc(auth.PasswordResetRequestPath, auth.createResetToken, "POST")
		router.HandleFunc(auth.passwordResetFormPath, auth.passwordResetForm, "GET")
		router.HandleFunc(auth.PasswordResetPath, auth.updatePassword, "POST")

		logrus.Infof(
			"password reset routes: token request -> %v, reset form -> %v, update password -> %v",
			auth.PasswordResetRequestPath,
			auth.passwordResetFormPath,
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
	reqToken := strings.Split(r.Header.Get("Authorization"), "Bearer")

	if len(reqToken) != 2 {
		return 0, "", "", errors.New("invalid token format")
	}

	token := strings.TrimSpace(reqToken[1])

	if token == "" {
		return 0, "", "", errors.New(fmt.Sprint("missing token field"))
	}

	claims, valid := auth.Token.ValidateAndExtract(token)
	if claims == nil || !valid {
		return 0, "", "", errors.New(fmt.Sprint("token validation failed: ", token))
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
		ConfirmRegistrationPath: auth.ServiceURL + accountConfirm,
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

	authUsecase := &user.AuthUsecase{
		Repository: auth.UserAccountRepository,
		Token:      auth.Token,
		Session: &user.Session{
			Cache: auth.CacheClient,
		},
	}

	var vars = auth.Router.Vars(r)

	account.ActivationToken = vars["token"]

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
		Repository:            auth.PasswordResetRepository,
		ConfirmResetTokenPath: auth.ServiceURL + passwordResetPath,
	}

	response := passwordResetUsecase.CreateResetToken(passwordReset.Email)

	w.Write(response.ToBytes())
}

// passwordResetForm API endpoint will show password reset form
func (auth *Auth) passwordResetForm(w http.ResponseWriter, r *http.Request) {
	if auth.PasswordResetCallback == nil {
		onError(errors.New("missing PasswordResetCallback implementation"), w)
		return
	}

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
	logrus.Error(err)
}

// onSuccess default callback
func onSuccess(payload []byte, w http.ResponseWriter) {
	w.Write(payload)
}

// passwordResetCallback default callback
func (auth *Auth) passwordResetCallback(w http.ResponseWriter, r *http.Request) {
	var vars = r.URL.Query()

	passwordReset := &user.PasswordReset{}
	passwordReset.Token = vars.Get("token")
	passwordReset.ServiceURL = auth.ServiceURL

	t, err := template.ParseFiles(passwordResetTemplate)
	if err != nil {
		w.Write([]byte("password reset template parse failed: " + err.Error()))
		return
	}

	if err := t.ExecuteTemplate(w, t.Name(), passwordReset); err != nil {
		w.Write([]byte("password reset template execute failed: " + err.Error()))
		return
	}
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

	if auth.PasswordResetCallback == nil {
		auth.PasswordResetCallback = auth.passwordResetCallback
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

	if strings.TrimSpace(auth.confirmRegistrationPath) == "" {
		auth.confirmRegistrationPath = accountConfirm + "{token}"
	}

	if strings.TrimSpace(auth.PasswordResetRequestPath) == "" {
		auth.PasswordResetRequestPath = passwordResetPath
	}

	if strings.TrimSpace(auth.passwordResetFormPath) == "" {
		auth.passwordResetFormPath = passwordResetPath
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
