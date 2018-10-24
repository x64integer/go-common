package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Channel for each ws connection
type Channel struct {
	ReadLock  sync.Mutex
	SendLock  sync.Mutex
	Conn      *websocket.Conn
	OnMessage func(in []byte)
	OnError   func(err error)
	Error     chan error
}

// NewChannel will init new Channel struct
func NewChannel(conn *websocket.Conn) *Channel {
	return &Channel{
		Conn:  conn,
		Error: make(chan error),
	}
}

// SendMessage is concurrent safe ws WriteMessage wrapper
func (ch *Channel) SendMessage(messageType int, data []byte) error {
	ch.SendLock.Lock()
	err := ch.Conn.WriteMessage(messageType, data)
	ch.SendLock.Unlock()

	return err
}

// ReadMessage is concurrent safe ws ReadMessage wrapper
func (ch *Channel) ReadMessage() (int, []byte, error) {
	ch.ReadLock.Lock()
	t, p, err := ch.Conn.ReadMessage()
	ch.ReadLock.Unlock()

	return t, p, err
}

// Read data from ws connection and pass it to OnMessage callback
func (ch *Channel) Read() {
	for {
		_, p, err := ch.ReadMessage()

		if err != nil {
			ch.Error <- err
			continue
		}

		ch.OnMessage(p)
	}
}

// HandleError from errors channel
func (ch *Channel) HandleError() {
	for {
		select {
		case err := <-ch.Error:
			ch.OnError(err)
		}
	}
}
