package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

// Authenticatable contract
// Such empty definition is used to satisfy reflection dependencies only
// It might change in the future, if such need arises
type Authenticatable interface{}

// Auth configuration
type Auth struct {
	RegisterPath string
	LoginPath    string
	LogoutPath   string
	Entity       Authenticatable
	OnError      func(error, http.ResponseWriter)
	OnSuccess    func([]byte, http.ResponseWriter)
}

// entityField is helper struct to hold information/data from extracted auth Entity (Authenticatable)
type entityField struct {
	Key         string
	Value       interface{}
	Type        interface{}
	ReflectType reflect.Type
}

// applyRoutes will setup auth routes (register, login, logout)
func (auth *Auth) applyRoutes(routeHandler RouteHandler) {
	if auth.OnError == nil {
		auth.OnError = onError
	}

	if auth.OnSuccess == nil {
		auth.OnSuccess = onSuccess
	}

	routeHandler.HandleFunc(auth.RegisterPath, func(w http.ResponseWriter, r *http.Request) {
		auth.handleFunc(w, r, func(fields []*entityField) ([]byte, error) {
			svc := &Service{}

			return svc.Register(fields)
		})
	})

	routeHandler.HandleFunc(auth.LoginPath, func(w http.ResponseWriter, r *http.Request) {
		auth.handleFunc(w, r, func(fields []*entityField) ([]byte, error) {
			svc := &Service{}

			return svc.Login(fields)
		})
	})

	routeHandler.HandleFunc(auth.LogoutPath, func(w http.ResponseWriter, r *http.Request) {
		auth.handleFunc(w, r, func(fields []*entityField) ([]byte, error) {
			svc := &Service{}

			return svc.Logout(fields)
		})
	})
}

// extractEntity is helper function to extract auth entity fields and tags
func (auth *Auth) extractEntity(entityToExtract interface{}) []*entityField {
	var entityType reflect.Type
	var entityValue reflect.Value
	var fields []*entityField
	entityKind := reflect.ValueOf(entityToExtract).Kind()

	if entityKind == reflect.Ptr {
		entityType = reflect.TypeOf(entityToExtract).Elem()
		entityValue = reflect.ValueOf(entityToExtract).Elem()
	} else {
		entityType = reflect.TypeOf(entityToExtract)
		entityValue = reflect.ValueOf(entityToExtract)
	}

	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		fieldKey := field.Tag.Get("auth")
		fieldValue := entityValue.Field(i)
		fieldType := field.Tag.Get("auth_type")

		fields = append(fields, &entityField{
			Key:         fieldKey,
			Value:       fieldValue,
			Type:        fieldType,
			ReflectType: field.Type,
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
	}

	fields := auth.extractEntity(auth.Entity)

	resp, err := callback(fields)
	if err != nil {
		auth.OnError(err, w)
		return
	}

	auth.OnSuccess(resp, w)
}

// onError default callback
func onError(err error, w http.ResponseWriter) {
	log.Println(err)
}

// onSuccess default callback
func onSuccess(payload []byte, w http.ResponseWriter) {
	w.Write(payload)
}
