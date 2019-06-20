package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/semirm-dev/go-common/api"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Server for websocket
type Server struct {
	EventHandler
	Config   *Config
	Channel  *Channel
	Upgrader websocket.Upgrader
}

// Run will create websocket Server and start listening for messages
func (server *Server) Run(done chan bool) {
	if server.Config == nil {
		log.Fatalln("nil Config struct for websocket server -> make sure valid Config is accessible to websocket server")
	}

	router := &api.MuxRouterAdapter{Router: mux.NewRouter()}

	router.HandleFunc(server.Config.Endpoint, func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			w.Write([]byte("failed to setup websocket upgrader"))
			return
		}

		ch := &Channel{
			Connection:   conn,
			EventHandler: server.EventHandler,
		}

		server.Channel = ch

		go ch.read()
	})

	go router.Listen(&api.Config{
		Host: server.Config.Host,
		Port: server.Config.Port,
	})

	<-done

	log.Println("reading stopped")
}

// SendText message to websocket channel
func (server *Server) SendText(msg []byte) error {
	return server.Channel.sendMessage(websocket.TextMessage, msg)
}

// SendBinary message to websocket channel
func (server *Server) SendBinary(msg []byte) error {
	return server.Channel.sendMessage(websocket.BinaryMessage, msg)
}
