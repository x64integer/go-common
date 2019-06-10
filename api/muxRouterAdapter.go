package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// MuxRouterAdapter will make sure mux.Router is compatible with Handler interface
type MuxRouterAdapter struct {
	*mux.Router
}

// Handle implements Handler.Handle()
func (muxRouterAdapter *MuxRouterAdapter) Handle(path string, handler http.Handler, methods ...string) {
	muxRouterAdapter.Router.Handle(path, handler).Methods(methods...)
}

// HandleFunc implements Handler.HandleFunc()
func (muxRouterAdapter *MuxRouterAdapter) HandleFunc(path string, f func(http.ResponseWriter, *http.Request), methods ...string) {
	muxRouterAdapter.Router.HandleFunc(path, f).Methods(methods...)
}
