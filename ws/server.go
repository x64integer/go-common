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
	MessageHandler
	Config   *Config
	Channel  *Channel
	Upgrader websocket.Upgrader
}

// Run will create websocket Server and start listening for messages
func (server *Server) Run(done chan bool) {
	if server.Config == nil || server.MessageHandler == nil {
		log.Fatalln("either Config or MessageHandler is missing")
	}

	router := &api.MuxRouterAdapter{Router: mux.NewRouter()}

	router.HandleFunc(server.Config.Endpoint, func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			w.Write([]byte("failed to setup websocket upgrader"))
			return
		}

		ch := &Channel{
			Connection:     conn,
			MessageHandler: server.MessageHandler,
		}

		server.Channel = ch

		go ch.read()
	}, "GET")

	go router.Listen(&api.Config{
		Host: server.Config.Host,
		Port: server.Config.Port,
	})

	<-done

	log.Println("server returned")
}

// SendText message to websocket channel
func (server *Server) SendText(msg []byte) error {
	return server.Channel.sendMessage(TextMessage, msg)
}

// SendBinary message to websocket channel
func (server *Server) SendBinary(msg []byte) error {
	return server.Channel.sendMessage(BinaryMessage, msg)
}
