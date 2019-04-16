package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// MuxRouterAdapter will make sure mux.Router is compatible with RouteHandler interface
type MuxRouterAdapter struct {
	*mux.Router
}

// Handle implements RouteHandler.Handle()
func (muxRouterAdapter *MuxRouterAdapter) Handle(path string, handler http.Handler, methods ...string) {
	muxRouterAdapter.Router.Handle(path, handler).Methods(methods...)
}

// HandleFunc implements RouteHandler.HandleFunc()
func (muxRouterAdapter *MuxRouterAdapter) HandleFunc(path string, f func(http.ResponseWriter, *http.Request), methods ...string) {
	muxRouterAdapter.Router.HandleFunc(path, f).Methods(methods...)
}
