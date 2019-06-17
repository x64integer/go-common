package mail

import (
	"errors"
	"sync"
)

// ErrMissingSender error
var ErrMissingSender = errors.New("missing Sender implementation")

// Sender for sending email
type Sender interface {
	Send(*Content) error
}

// Client is used to safely send email
type Client struct {
	Sender
	SendLock sync.Mutex
}

// Send email
func (client *Client) Send(content *Content) error {
	client.SendLock.Lock()
	defer client.SendLock.Unlock()

	if client.Sender == nil {
		return ErrMissingSender
	}

	return client.Sender.Send(content)
}
