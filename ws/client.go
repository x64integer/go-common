package ws

import (
	"errors"

	"github.com/gorilla/websocket"
)

// Client for websocket
type Client struct {
	Config      *Config
	Channel     *Channel
	OnMessage   func(in []byte)
	OnError     func(err error)
	OnConnClose func(code int, msg string)
}

// Setup will create websocket Client and start listening for messages
func (client *Client) Setup() error {
	if client.Config == nil {
		return errors.New("nil Config struct for ws Client -> make sure valid Config is accessible to ws Client")
	}

	c, _, err := websocket.DefaultDialer.Dial(client.Config.WSURL, nil)
	if err != nil {
		return err
	}

	ch := &Channel{
		Conn:        c,
		OnMessage:   client.OnMessage,
		OnError:     client.OnError,
		OnConnClose: client.OnConnClose,
	}

	client.Channel = ch

	if client.OnMessage != nil {
		go ch.Read()
	}

	return nil
}

// SendText message
func (client *Client) SendText(msg []byte) error {
	return client.Channel.SendMessage(websocket.TextMessage, msg)
}

// SendBinary message
func (client *Client) SendBinary(msg []byte) error {
	return client.Channel.SendMessage(websocket.BinaryMessage, msg)
}
