package mail

import "sync"

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

	return client.Sender.Send(content)
}
