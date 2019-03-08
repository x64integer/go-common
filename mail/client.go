package mail

import "sync"

// Sender for sending email
type Sender interface {
	Send(*Content) error
}

// Client ...
type Client struct {
	Sender
	SendLock sync.Mutex
}

// Send email
func (client *Client) Send(content *Content) error {
	client.SendLock.Lock()
	client.SendLock.Unlock()

	return client.Sender.Send(content)
}
