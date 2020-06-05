package ws

import (
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

const (
	// TextMessage type
	TextMessage = websocket.TextMessage
	// BinaryMessage type
	BinaryMessage = websocket.BinaryMessage
)

// MessageHandler for websocket connection
type MessageHandler interface {
	// OnMessage callback
	// func(messageType, content) is used for reply back
	OnMessage([]byte, func(int, []byte) error)
	// OnError callback
	OnError(error)
}

// Channel for websocket connection
type Channel struct {
	MessageHandler
	Connection *websocket.Conn
	ReadLock   sync.Mutex
	SendLock   sync.Mutex
}

// read data from websocket channel
// Concurrent safe wrapper for Connection.ReadMessage()
func (ch *Channel) read() {
	defer func() {
		logrus.Warn("websocket connection stopped reading messages")

		if err := ch.close(); err != nil {
			logrus.Error("failed to close websocket connection: ", err)
			return
		}

		logrus.Warn("websocket connection closed")
	}()

	logrus.Info("listening for messages")

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

// close websocket connection
func (ch *Channel) close() error {
	return ch.Connection.Close()
}
