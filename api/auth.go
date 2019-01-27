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
	Registrable
	Loginable
}

// applyRoutes will setup register and login routes
func (auth *Auth) applyRoutes(routeHandler RouteHandler) {
	routeHandler.HandleFunc(auth.RegisterPath, func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("failed to read register request body: ", err)
			return
		}

		entityToExtract := auth.Entity
		if auth.Registrable != nil {
			entityToExtract = auth.Registrable
		}

		entity := auth.extractEntity(entityToExtract)

		for k, v := range entity {
			log.Printf("Field: %v, Tag: %v, Type: %v", k, v.Tag, v.Type)
		}

		w.Write(b)
	})

	routeHandler.HandleFunc(auth.LoginPath, func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("failed to read login request body: ", err)
			return
		}

		entityToExtract := auth.Entity
		if auth.Loginable != nil {
			entityToExtract = auth.Loginable
		}

		entity := auth.extractEntity(entityToExtract)

		for k, v := range entity {
			log.Printf("Field: %v, Tag: %v, Type: %v", k, v.Tag, v.Type)
		}

		w.Write(b)
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
