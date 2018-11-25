package ws

import (
	"errors"

	"github.com/gorilla/websocket"
)

// Connection for ws
type Connection struct {
	Config    *Config
	Channel   *Channel
	OnMessage func(in []byte)
	OnError   func(err error)
}

// Setup will create ws connection and start listening for messages
func (conn *Connection) Setup() error {
	if conn.Config == nil {
		return errors.New("nil Config struct for ws connection -> make sure valid Config is accessible to ws connection")
	}

	c, _, err := websocket.DefaultDialer.Dial(conn.Config.WSURL, nil)
	if err != nil {
		return err
	}

	ch := &Channel{
		Conn:      c,
		OnMessage: conn.OnMessage,
		OnError:   conn.OnError,
	}

	conn.Channel = ch

	if conn.OnMessage != nil {
		go ch.Read()
	}

	if conn.OnError != nil {
		go ch.HandleError()
	}

	return nil
}

// SendText message
func (conn *Connection) SendText(msg []byte) error {
	return conn.Channel.SendMessage(websocket.TextMessage, msg)
}

// SendBinary message
func (conn *Connection) SendBinary(msg []byte) error {
	return conn.Channel.SendMessage(websocket.BinaryMessage, msg)
}
