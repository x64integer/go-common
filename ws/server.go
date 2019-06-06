package ws

import (
	"log"
	"net/http"

	"github.com/semirm-dev/go-common/api"

	"github.com/gorilla/websocket"
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

	r := api.NewRouter(&api.Config{
		Host: server.Config.Host,
		Port: server.Config.Port,
	})

	r.HandleFunc(server.Config.Endpoint, func(w http.ResponseWriter, r *http.Request) {
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

	go r.Listen()

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
