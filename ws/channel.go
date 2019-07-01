package ws

import (
	"log"
	"sync"
)

const (
	// TextMessage type
	TextMessage = 1
	// BinaryMessage type
	BinaryMessage = 2
)

// MessageHandler for websocket connection
type MessageHandler interface {
	// OnMessage callback
	// func(messageType, content) is used for reply back
	OnMessage([]byte, func(int, []byte) error)
	// OnError callback
	OnError(error)
}

// Connection for websocket
type Connection interface {
	// ReadMessage from websocket connection
	ReadMessage() (int, []byte, error)
	// WriteMessage to websocket connection
	WriteMessage(int, []byte) error
	// Close websocket connection
	Close() error
}

// Channel for websocket connection
type Channel struct {
	MessageHandler
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
	defer func() {
		if err := ch.Connection.Close(); err != nil {
			log.Println("failed to close websocket connection: ", err)
			return
		}

		log.Println("websocket connection closed")
	}()

	log.Println("listening for messages")

	for {
		_, msg, err := ch.readMessage()

		if err != nil {
			ch.MessageHandler.OnError(err)

			break
		}

		ch.MessageHandler.OnMessage(msg, func(t int, b []byte) error {
			return ch.sendMessage(t, b)
		})
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
