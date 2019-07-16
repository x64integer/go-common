package ws

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/semirm-dev/go-common/api"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Server for websocket
type Server struct {
	MessageHandler
	Host     string
	Port     string
	Endpoint string
	Channel  *Channel
}

// Run will create websocket Server and start listening for messages
func (server *Server) Run(done chan bool) {
	if server.MessageHandler == nil {
		logrus.Fatal("MessageHandler is missing")
	}

	router := &api.MuxRouterAdapter{Router: mux.NewRouter()}

	router.HandleFunc(server.Endpoint, func(w http.ResponseWriter, r *http.Request) {
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
		Host: server.Host,
		Port: server.Port,
	})

	<-done

	logrus.Warn("server returned")
}

// SendText message to websocket channel
func (server *Server) SendText(msg []byte) error {
	return server.Channel.sendMessage(TextMessage, msg)
}

// SendBinary message to websocket channel
func (server *Server) SendBinary(msg []byte) error {
	return server.Channel.sendMessage(BinaryMessage, msg)
}
