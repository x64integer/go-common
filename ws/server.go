package ws

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
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

	router := mux.NewRouter()

	router.HandleFunc(server.Endpoint, func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			_, err := w.Write([]byte("failed to setup websocket upgrader"))
			if err != nil {
				logrus.Error("failed to send response from upgrader: ", err.Error())
			}

			return
		}

		ch := &Channel{
			Connection:     conn,
			MessageHandler: server.MessageHandler,
		}

		server.Channel = ch

		go ch.read()
	})

	go func() {
		if err := http.ListenAndServe(server.Host+":"+server.Port, router); err != nil {
			logrus.Fatal("failed to start http server: ", err.Error())
		}
	}()

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
