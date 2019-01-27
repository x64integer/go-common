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

// Logoutable contract is used in case we want different entity for logout
type Logoutable interface{}

// Auth configuration
type Auth struct {
	RegisterPath string
	LoginPath    string
	LogoutPath   string
	Entity       Authenticatable
	// Optional properties for customization
	*Registration
	*Login
	*Logout
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

// Logout customization
type Logout struct {
	Path      string
	Entity    Loginable
	OnError   func(error, http.ResponseWriter)
	OnSuccess func([]byte, http.ResponseWriter)
}

// applyRoutes will setup auth routes (register, login, logout)
func (auth *Auth) applyRoutes(routeHandler RouteHandler) {
	registerPath, registerEntity, onRegisterError, onRegisterSuccess := auth.mapRegistration()
	loginPath, loginEntity, onLoginError, onLoginSuccess := auth.mapLogin()
	logoutPath, logoutEntity, onLogoutError, onLogoutSuccess := auth.mapLogout()

	routeHandler.HandleFunc(registerPath, func(w http.ResponseWriter, r *http.Request) {
		auth.handleFunc(w, r, registerEntity, onRegisterError, onRegisterSuccess, func(payload []byte) {
			// register logic
		})
	})

	routeHandler.HandleFunc(loginPath, func(w http.ResponseWriter, r *http.Request) {
		auth.handleFunc(w, r, loginEntity, onLoginError, onLoginSuccess, func(payload []byte) {
			// login logic
		})
	})

	routeHandler.HandleFunc(logoutPath, func(w http.ResponseWriter, r *http.Request) {
		auth.handleFunc(w, r, logoutEntity, onLogoutError, onLogoutSuccess, func(payload []byte) {
			// logout logic
		})
	})
}

// authEntity is helper struct to hold information/data from extracted auth Entity (Authenticatable, Registrable, Loginable, Logoutable)
type authEntity struct {
	Field string
	Tag   string
	Type  interface{}
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

// mapRegistration is helper function to map registration data structures
// Initiate default values
// Override with values defined in auth.Registration struct
func (auth *Auth) mapRegistration() (string, Registrable, func(error, http.ResponseWriter), func([]byte, http.ResponseWriter)) {
	registerPath := auth.RegisterPath
	registerEntity := auth.Entity
	onRegisterError := func(err error, w http.ResponseWriter) {
		log.Println("registration failed: ", err)
	}
	onRegisterSuccess := func(payload []byte, w http.ResponseWriter) {
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

	return registerPath, registerEntity, onRegisterError, onRegisterSuccess
}

// mapLogin is helper function to map login data structures
// Initiate default values
// Override with values defined in auth.Login struct
func (auth *Auth) mapLogin() (string, Loginable, func(error, http.ResponseWriter), func([]byte, http.ResponseWriter)) {
	loginPath := auth.LoginPath
	loginEntity := auth.Entity
	onLoginError := func(err error, w http.ResponseWriter) {
		log.Println("login failed: ", err)
	}
	onLoginSuccess := func(payload []byte, w http.ResponseWriter) {
		w.Write(payload)
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

	return loginPath, loginEntity, onLoginError, onLoginSuccess
}

// mapLogout is helper function to map logout data structures
// Initiate default values
// Override with values defined in auth.Logout struct
func (auth *Auth) mapLogout() (string, Logoutable, func(error, http.ResponseWriter), func([]byte, http.ResponseWriter)) {
	logoutPath := auth.LogoutPath
	logoutEntity := auth.Entity
	onLogoutError := func(err error, w http.ResponseWriter) {
		log.Println("logout failed: ", err)
	}
	onLogoutSuccess := func(payload []byte, w http.ResponseWriter) {
		w.Write(payload)
	}

	if auth.Logout != nil {
		if auth.Logout.Path != "" {
			logoutPath = auth.Logout.Path
		}

		if auth.Logout.Entity != nil {
			logoutEntity = auth.Logout.Entity
		}

		if auth.Logout.OnError != nil {
			onLogoutError = auth.Logout.OnError
		}

		if auth.Logout.OnSuccess != nil {
			onLogoutSuccess = auth.Logout.OnSuccess
		}
	}

	return logoutPath, logoutEntity, onLogoutError, onLogoutSuccess
}

// handleFunc is helper function to setup route
func (auth *Auth) handleFunc(
	w http.ResponseWriter,
	r *http.Request,
	entityToExtract interface{},
	onError func(error, http.ResponseWriter),
	onSuccess func([]byte, http.ResponseWriter),
	callback func(payload []byte),
) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		onError(err, w)
		return
	}

	entity := auth.extractEntity(entityToExtract)

	for k, v := range entity {
		log.Printf("Field: %v, Tag: %v, Type: %v", k, v.Tag, v.Type)
	}

	callback(b)

	onSuccess(b, w)
}
