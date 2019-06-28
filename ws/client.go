package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

// Client for websocket
type Client struct {
	EventHandler
	Config         *Config
	Channel        *Channel
	DisabledReader bool
}

// Connect will create websocket Client and start listening for messages
func (client *Client) Connect(done chan bool, ready chan bool) {
	if client.Config == nil {
		log.Fatalln("nil Config struct for websocket Client -> make sure valid Config is accessible to websocket Client")
	}

	conn := client.connect()

	defer func() {
		log.Println("connection closed")
		defer conn.Close()
	}()

	ch := &Channel{
		Connection:   conn,
		EventHandler: client.EventHandler,
	}

	client.Channel = ch

	if !client.DisabledReader {
		go ch.read()
	}

	ready <- true

	log.Printf("\nconnected to %s\n", client.Config.WSURL)

	<-done
}

// SendText message to websocket channel
func (client *Client) SendText(msg []byte) error {
	return client.Channel.sendMessage(websocket.TextMessage, msg)
}

// SendBinary message to websocket channel
func (client *Client) SendBinary(msg []byte) error {
	return client.Channel.sendMessage(websocket.BinaryMessage, msg)
}

// connect is helper function to create gorilla websocket connection
func (client *Client) connect() *websocket.Conn {
	conn, _, err := websocket.DefaultDialer.Dial(client.Config.WSURL, nil)
	if err != nil {
		log.Fatalln("websocket dialer failed: ", err)
	}

	return conn
}
