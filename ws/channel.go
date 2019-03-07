package ws

import (
	"log"
	"sync"
)

// EventHandler for websocket connection
type EventHandler interface {
	// OnMessage event from websocket connection
	OnMessage(in []byte)
	// OnError event from websocket connection
	OnError(err error)
}

// Connection for websocket
type Connection interface {
	// ReadMessage from websocket connection
	ReadMessage() (int, []byte, error)
	// WriteMessage to websocket connection
	WriteMessage(int, []byte) error
}

// Channel for websocket connection
type Channel struct {
	EventHandler
	Connection
	ReadLock sync.Mutex
	SendLock sync.Mutex
	Error    chan error
}

// ConnectionClosed error type
type ConnectionClosed struct {
	Code    int
	Message string
}

// Error implements error interface
func (connClosed *ConnectionClosed) Error() string {
	return "connection closed"
}

// read data from websocket channel
// Concurrent safe wrapper for Connection.ReadMessage()
func (ch *Channel) read() {
	log.Println("listening for messages")

	for {
		_, msg, err := ch.readMessage()

		if err != nil {
			ch.EventHandler.OnError(err)

			continue
		}

		ch.EventHandler.OnMessage(msg)
	}
}

// sendMessage to websocket channel
func (ch *Channel) sendMessage(messageType int, data []byte) error {
	ch.SendLock.Lock()
	err := ch.Connection.WriteMessage(messageType, data)
	ch.SendLock.Unlock()

	return err
}

// readMessage from websocket channel
func (ch *Channel) readMessage() (int, []byte, error) {
	ch.ReadLock.Lock()
	t, p, err := ch.Connection.ReadMessage()
	ch.ReadLock.Unlock()

	return t, p, err
}
