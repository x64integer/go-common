package mail

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/x64integer/go-common/util"
)

var (
	// From - mail account to send mails from (admin mail or something)
	From = util.Env("MAIL_SERVICE_FROM", "")
	// Password for sending mail account
	Password = util.Env("MAIL_SERVICE_FROM_PASSWORD", "")
)

// Entity holds information about sending email
type Entity struct {
	To         []string
	Cc         []string
	Bcc        []string
	Subject    string
	Content    []byte
	Attachment []byte
}

// SMTPServer for mail
type SMTPServer struct {
	Host      string
	Port      string
	TLSConfig *tls.Config
}

// NewSMTPServer will initialize SMTP server
func NewSMTPServer() *SMTPServer {
	smtpServer := &SMTPServer{
		Host: util.Env("MAIL_SERVICE_SMTP_HOST", "smtp.gmail.com"),
		Port: util.Env("MAIL_SERVICE_SMTP_PORT", "465"),
	}

	smtpServer.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer.Host,
	}

	return smtpServer
}

// Send mail
func (entity *Entity) Send() error {
	client, err := entity.client()
	if err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(entity.body())
	if err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return nil
	}

	client.Quit()

	return nil
}

// client with required mail configurations
func (entity *Entity) client() (*smtp.Client, error) {
	smtpServer := NewSMTPServer()

	conn, err := tls.Dial("tcp", smtpServer.Host+":"+smtpServer.Port, smtpServer.TLSConfig)
	if err != nil {
		return nil, err
	}

	client, err := smtp.NewClient(conn, smtpServer.Host)
	if err != nil {
		return nil, err
	}

	auth := smtp.PlainAuth("", From, Password, smtpServer.Host)

	if err := client.Auth(auth); err != nil {
		return nil, err
	}

	if err := client.Mail(From); err != nil {
		return nil, err
	}

	receivers := append(entity.To, entity.Cc...)
	receivers = append(receivers, entity.Bcc...)

	for _, k := range receivers {
		if err = client.Rcpt(k); err != nil {
			continue
		}
	}

	return client, nil
}

// body mail to be sent
func (entity *Entity) body() []byte {
	header := ""

	header += fmt.Sprintf("From: %s\r\n", From)

	if len(entity.To) > 0 {
		header += fmt.Sprintf("To: %s\r\n", strings.Join(entity.To, ";"))
	}

	if len(entity.Cc) > 0 {
		header += fmt.Sprintf("Cc: %s\r\n", strings.Join(entity.Cc, ";"))
	}

	header += fmt.Sprintf("Subject: %s\r\n", entity.Subject)
	header += "\r\n" + string(entity.Content)

	return []byte(header)
}
