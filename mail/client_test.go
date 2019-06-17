package mail_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/semirm-dev/go-common/mail"
)

var content = &mail.Content{
	From: "mail-sender@mail.com",
	To: []string{
		"mail-receiver-1@mail.com",
		"mail-receiver-2@mail.com",
	},
	Cc: []string{
		"mail-receiver-cc-1@mail.com",
	},
	Bcc: []string{
		"mail-receiver-1-bcc@mail.com",
	},
	Subject: "mail-subject",
	Body:    []byte("mail body"),
}

type Mock struct{}

// Send mail implements Sender.Send
func (mock *Mock) Send(content *mail.Content) error {
	return nil
}

func TestClientSend(t *testing.T) {
	client := &mail.Client{
		Sender: &Mock{},
	}

	err := client.Send(content)

	assert.NoError(t, err)
}

func TestClientSendNilSender(t *testing.T) {
	client := &mail.Client{}

	err := client.Send(content)

	assert.Error(t, err)
	assert.Equal(t, mail.ErrMissingSender, err)
}
