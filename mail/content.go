package mail

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
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

// Construct content for email
func (content *Content) Construct() []byte {
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

	tpl := template.Must(template.New("emailTemplate").Parse(header))
	if err := tpl.Execute(buffer, &content); err != nil {
		logrus.Error("failed to execute template: ", err.Error())
		return nil
	}

	return buffer.Bytes()
}
