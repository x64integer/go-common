package ws

import "github.com/gorilla/websocket"

// Connection for ws
type Connection struct {
	Config    *Config
	OnMessage func(in []byte)
	OnError   func(err error)
	Channel   *Channel
	KeepAlive bool
}

// Setup will create ws connection and start listening for messages
func (conn *Connection) Setup() error {
	c, _, err := websocket.DefaultDialer.Dial(conn.Config.WSURL, nil)
	if err != nil {
		return err
	}

	if !conn.KeepAlive {
		defer c.Close()
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

	go ch.HandleError()

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
