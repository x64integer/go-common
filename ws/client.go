package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

// Client for websocket
type Client struct {
	MessageHandler
	Config         *Config
	Channel        *Channel
	DisabledReader bool
}

// Connect will create websocket Client and start listening for messages
func (client *Client) Connect(done chan bool, ready chan bool) {
	if client.Config == nil || client.MessageHandler == nil {
		log.Fatalln("either Config or MessageHandler is missing")
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

	log.Println("connected to: ", client.Config.WSURL)

	<-done

	log.Println("client returned")
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
		log.Fatalln("websocket dialer failed: ", err)
	}

	return conn
}
