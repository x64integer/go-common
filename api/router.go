package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

// RouteHandler decouples direct dependency on *mux.Route
// Add new functions as per need
type RouteHandler interface {
	Handle(string, http.Handler)
	HandleFunc(string, func(http.ResponseWriter, *http.Request))
}

// Router for api
type Router struct {
	*Config
	*http.Server
}

// Config for router
type Config struct {
	Host        string
	Port        string
	MapRoutes   func(RouteHandler)
	WaitTimeout time.Duration
	*Auth
}

// NewRouter will initialize new Router
func NewRouter(config *Config) *Router {
	gRouter := mux.NewRouter()

	routeHandler := &MuxRouterAdapter{Router: gRouter}

	if config.Auth != nil {
		config.Auth.applyRoutes(routeHandler)
		log.Printf("registered auth routes: register -> %v, login -> %v", config.Auth.RegisterPath, config.Auth.LoginPath)
	}

	config.MapRoutes(routeHandler)

	srv := &http.Server{
		Addr: config.Host + ":" + config.Port,
		// NOTE: these values could be passed through Config
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      gRouter,
	}

	return &Router{
		Config: config,
		Server: srv,
	}
}

// Listen http server on our router
func (r *Router) Listen() {
	go func() {
		log.Printf("listening on %v:%v\n", r.Config.Host, r.Config.Port)
		if err := r.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// NOTE: as per gorilla/mux preferences
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), r.Config.WaitTimeout)
	defer cancel()

	r.Shutdown(ctx)

	log.Println("shutting down http server")
	os.Exit(0)
}
