package ws

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/x64integer/go-common/api"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Server for websocket
type Server struct {
	Config      *Config
	Channel     *Channel
	OnMessage   func(in []byte)
	OnError     func(err error)
	OnConnClose func(code int, msg string)
	Upgrader    websocket.Upgrader
}

// Run will create websocket Server and start listening for messages
func (server *Server) Run(done chan bool) {
	if server.Config == nil {
		log.Fatalln("nil Config struct for ws server -> make sure valid Config is accessible to ws server")
	}

	r := api.NewRouter(&api.Config{
		Host:        server.Config.Host,
		Port:        server.Config.Port,
		WaitTimeout: time.Second * 15,
		MapRoutes: func(r *mux.Router) {
			r.HandleFunc(server.Config.Endpoint, func(w http.ResponseWriter, r *http.Request) {
				c, err := upgrader.Upgrade(w, r, nil)
				if err != nil {
					w.Write([]byte("failed to setup websocket upgrader"))
					return
				}

				ch := &Channel{
					Conn:        c,
					OnMessage:   server.OnMessage,
					OnError:     server.OnError,
					OnConnClose: server.OnConnClose,
				}

				server.Channel = ch

				go ch.Read()
			})
		},
	})

	go r.Listen()

	<-done
}
