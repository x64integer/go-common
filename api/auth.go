package api

import (
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

// Authenticatable contract is used for both register and login
type Authenticatable interface{}

// Registrable contract is used in case we want different entity for registration
type Registrable interface{}

// Loginable contract is used in case we want different entity for login
type Loginable interface{}

// Auth configuration
type Auth struct {
	RegisterPath string
	LoginPath    string
	Entity       Authenticatable
	// Optional properties for customization
	*Registration
	*Login
}

// Registration customization
type Registration struct {
	Path      string
	Entity    Registrable
	OnError   func(error, http.ResponseWriter)
	OnSuccess func([]byte, http.ResponseWriter)
}

// Login customization
type Login struct {
	Path      string
	Entity    Loginable
	OnError   func(error, http.ResponseWriter)
	OnSuccess func([]byte, http.ResponseWriter)
}

// applyRoutes will setup register and login routes
func (auth *Auth) applyRoutes(routeHandler RouteHandler) {
	registerPath, loginPath := auth.RegisterPath, auth.LoginPath
	registerEntity, loginEntity := auth.Entity, auth.Entity
	onRegisterError, onLoginError := func(err error, w http.ResponseWriter) {
		log.Println("registration failed: ", err)
	}, func(err error, w http.ResponseWriter) {
		log.Println("login failed: ", err)
	}
	onRegisterSuccess, onLoginSuccess := func(payload []byte, w http.ResponseWriter) {
		w.Write(payload)
	}, func(payload []byte, w http.ResponseWriter) {
		w.Write(payload)
	}

	if auth.Registration != nil {
		if auth.Registration.Path != "" {
			registerPath = auth.Registration.Path
		}

		if auth.Registration.Entity != nil {
			registerEntity = auth.Registration.Entity
		}

		if auth.Registration.OnError != nil {
			onRegisterError = auth.Registration.OnError
		}

		if auth.Registration.OnSuccess != nil {
			onRegisterSuccess = auth.Registration.OnSuccess
		}
	}

	if auth.Login != nil {
		if auth.Login.Path != "" {
			loginPath = auth.Login.Path
		}

		if auth.Login.Entity != nil {
			loginEntity = auth.Login.Entity
		}

		if auth.Login.OnError != nil {
			onLoginError = auth.Login.OnError
		}

		if auth.Login.OnSuccess != nil {
			onLoginSuccess = auth.Login.OnSuccess
		}
	}

	routeHandler.HandleFunc(registerPath, func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			onRegisterError(err, w)
			return
		}

		entity := auth.extractEntity(registerEntity)

		for k, v := range entity {
			log.Printf("Field: %v, Tag: %v, Type: %v", k, v.Tag, v.Type)
		}

		onRegisterSuccess(b, w)
	})

	routeHandler.HandleFunc(loginPath, func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			onLoginError(err, w)
			return
		}

		entity := auth.extractEntity(loginEntity)

		for k, v := range entity {
			log.Printf("Field: %v, Tag: %v, Type: %v", k, v.Tag, v.Type)
		}

		onLoginSuccess(b, w)
	})
}

// extractEntity is helper function to extract auth entity fields and tags
func (auth *Auth) extractEntity(entityToExtract interface{}) map[string]*authEntity {
	var entityType reflect.Type
	entityKind := reflect.ValueOf(entityToExtract).Kind()
	entityExtracted := make(map[string]*authEntity)

	if entityKind == reflect.Ptr {
		entityType = reflect.TypeOf(entityToExtract).Elem()
	} else {
		entityType = reflect.TypeOf(entityToExtract)
	}

	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		authTag := field.Tag.Get("auth")
		authType := field.Tag.Get("auth_type")

		entityExtracted[field.Name] = &authEntity{
			Field: field.Name,
			Tag:   authTag,
			Type:  authType,
		}
	}

	return entityExtracted
}

// authEntity is helper struct to hold information/data from extracted auth Entity (Authenticatable, Registrable, Loginable)
type authEntity struct {
	Field string
	Tag   string
	Type  interface{}
}
