package api

import (
	"net/http"
)

// AuthMiddleware is responsible to protect api routes to registered users only
var AuthMiddleware = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// NOTE: add logic for authentication
		// get auth token and validate it
		// we could also use authO.com

		next.ServeHTTP(w, r)
	})
}
