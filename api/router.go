package api

import (
	"net/http"
)

// Router will listen and serve http routes
// Add new functions as per need
type Router interface {
	Listen(*Config)
	Handle(string, http.Handler, ...string)
	HandleFunc(string, func(http.ResponseWriter, *http.Request), ...string)
}

// Config for router
type Config struct {
	Host string
	Port string
}
