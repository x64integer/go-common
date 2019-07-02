package ws

import (
	"github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

// Client for websocket
type Client struct {
	Config *Config
	MessageHandler
	Channel        *Channel
	DisabledReader bool
}

// Connect will create websocket Client and start listening for messages
func (client *Client) Connect(done chan bool, ready chan bool) {
	if client.Config == nil || client.MessageHandler == nil {
		logrus.Fatal("either Config or MessageHandler is missing")
	}

	conn := client.connection()

	ch := &Channel{
		Connection:     conn,
		MessageHandler: client.MessageHandler,
	}

	client.Channel = ch

	if !client.DisabledReader {
		go ch.read()
	}

	ready <- true

	logrus.Info("connected to: ", client.Config.WSURL)

	<-done

	logrus.Warn("client returned")
}

// SendText message to websocket channel
func (client *Client) SendText(msg []byte) error {
	return client.Channel.sendMessage(TextMessage, msg)
}

// SendBinary message to websocket channel
func (client *Client) SendBinary(msg []byte) error {
	return client.Channel.sendMessage(BinaryMessage, msg)
}

// connection is helper function to create gorilla websocket connection
func (client *Client) connection() *websocket.Conn {
	conn, _, err := websocket.DefaultDialer.Dial(client.Config.WSURL, nil)
	if err != nil {
		logrus.Fatal("websocket dialer failed: ", err)
	}

	return conn
}
