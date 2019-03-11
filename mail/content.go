package mail

import (
	"fmt"
	"strings"
)

// Content holds information about email
type Content struct {
	From       string
	To         []string
	Cc         []string
	Bcc        []string
	Subject    string
	Body       []byte
	Attachment []byte
}

// construct content for email
func (content *Content) construct() []byte {
	header := ""

	header += fmt.Sprintf("From: %s\r\n", content.From)

	if len(content.To) > 0 {
		header += fmt.Sprintf("To: %s\r\n", strings.Join(content.To, ";"))
	}

	if len(content.Cc) > 0 {
		header += fmt.Sprintf("Cc: %s\r\n", strings.Join(content.Cc, ";"))
	}

	header += fmt.Sprintf("Subject: %s\r\n", content.Subject)
	header += "\r\n" + string(content.Body)

	return []byte(header)
}
