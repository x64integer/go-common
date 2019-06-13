package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

// Handler decouples direct dependency on *mux.Route.
// Add new functions as per need
type Handler interface {
	Handle(string, http.Handler, ...string)
	HandleFunc(string, func(http.ResponseWriter, *http.Request), ...string)
}

// Router for api
type Router struct {
	*Config
	*http.Server
	Handler
}

// Config for router
type Config struct {
	Host string
	Port string
	*Auth
}

// NewRouter will initialize new Router
func NewRouter(config *Config) *Router {
	muxRouter := mux.NewRouter()

	handler := &MuxRouterAdapter{Router: muxRouter}

	addr := config.Host + ":" + config.Port

	if config.Auth != nil {
		config.Auth.serviceURL = addr
		config.Auth.apply(handler)
	}

	srv := &http.Server{
		Addr: addr,
		// NOTE: these values could be passed through Config
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler,
	}

	return &Router{
		Config:  config,
		Server:  srv,
		Handler: handler,
	}
}

// Listen http server on our router
func (r *Router) Listen() {
	go func() {
		logrus.Infof("listening on %v:%v\n", r.Config.Host, r.Config.Port)
		if err := r.Server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// NOTE: as per gorilla/mux preferences
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	r.Shutdown(ctx)

	logrus.Warn("shutting down http server")
	os.Exit(0)
}
