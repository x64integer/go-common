package api

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// MuxRouterAdapter implements Router interface
type MuxRouterAdapter struct {
	*mux.Router
}

// Listen implements Router.Listen()
func (adapter *MuxRouterAdapter) Listen(config *Config) {
	server := &http.Server{
		Addr: config.Host + ":" + config.Port,
		// NOTE: these values could be passed through Config
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      adapter,
	}

	go func() {
		logrus.Infof("listening on %v\n", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c

	logrus.Warn("shutting down http server")
	os.Exit(0)
}

// Handle implements Router.Handle()
func (adapter *MuxRouterAdapter) Handle(path string, handler http.Handler, methods ...string) {
	adapter.Router.Handle(path, handler).Methods(methods...)
}

// HandleFunc implements Router.HandleFunc()
func (adapter *MuxRouterAdapter) HandleFunc(path string, f func(http.ResponseWriter, *http.Request), methods ...string) {
	adapter.Router.HandleFunc(path, f).Methods(methods...)
}
