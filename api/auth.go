package api

import (
	"encoding/json"
	"fmt"
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
	*Service
}

// entityField is helper struct to hold information/data from extracted auth Entity (Authenticatable)
type entityField struct {
	AuthKey   string
	AuthValue interface{}
	AuthType  interface{}
	AuthTable string
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
			return auth.Service.Register(fields)
		})
	})

	routeHandler.HandleFunc(auth.LoginPath, func(w http.ResponseWriter, r *http.Request) {
		auth.handleFunc(w, r, func(fields []*entityField) ([]byte, error) {
			return auth.Service.Login(fields)
		})
	})

	routeHandler.HandleFunc(auth.LogoutPath, func(w http.ResponseWriter, r *http.Request) {
		auth.handleFunc(w, r, func(fields []*entityField) ([]byte, error) {
			return auth.Service.Logout(fields)
		})
	})
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
			AuthKey:   fieldKey,
			AuthValue: fieldValue,
			AuthType:  fieldType,
			AuthTable: entityValue.Type().Name(),
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

// onError default callback
func onError(err error, w http.ResponseWriter) {
	log.Println(err)
}

// onSuccess default callback
func onSuccess(payload []byte, w http.ResponseWriter) {
	w.Write(payload)
}
