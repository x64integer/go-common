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
// We can still call Handle() on RouteHandler without *mux.Route in its return definition
// Just a little trick to get rid of *mux.Route dependency
// Add new functions as per need
type RouteHandler interface {
	Handle(path string, handler http.Handler) *mux.Route
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route
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
}

// NewRouter will initialize new Router
func NewRouter(conf *Config) *Router {
	gRouter := mux.NewRouter()

	conf.MapRoutes(gRouter)

	srv := &http.Server{
		Addr: conf.Host + ":" + conf.Port,
		// NOTE: these values could be passed through Config
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      gRouter,
	}

	return &Router{
		Config: conf,
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
