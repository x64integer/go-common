package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

// Content holds information about email
type Content struct {
	From        string
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	Body        []byte
	Attachment  []byte
	ContentType string
}

// construct content for email
func (content *Content) construct() []byte {
	if content.ContentType == "" {
		content.ContentType = "text/html"
	}

	header := ""

	header += fmt.Sprintf("From: %s\r\n", content.From)

	if len(content.To) > 0 {
		header += fmt.Sprintf("To: %s\r\n", strings.Join(content.To, ";"))
	}

	if len(content.Cc) > 0 {
		header += fmt.Sprintf("Cc: %s\r\n", strings.Join(content.Cc, ";"))
	}

	header += fmt.Sprintf("Subject: %s\r\n", content.Subject)

	header += fmt.Sprint("Content-Type: " + content.ContentType + "; charset=\"UTF-8\"\r\n")

	header += "\r\n" + string(content.Body)

	buffer := new(bytes.Buffer)

	template := template.Must(template.New("emailTemplate").Parse(header))
	template.Execute(buffer, &content)

	return buffer.Bytes()
}
