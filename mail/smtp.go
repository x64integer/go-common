package mail

import (
	"crypto/tls"
	"log"
	"net/smtp"

	"github.com/semirm-dev/go-common/util"
)

// SMTP for mail
type SMTP struct {
	From      string
	Password  string
	Host      string
	Port      string
	Client    *smtp.Client
	TLSConfig *tls.Config
}

// DefaultSMTP will initialize SMTP server with default values
func DefaultSMTP() *SMTP {
	smtpServer := &SMTP{
		From:     util.Env("MAIL_FROM", ""),
		Password: util.Env("MAIL_FROM_PASSWORD", ""),
		Host:     util.Env("MAIL_SMTP_HOST", "smtp.gmail.com"),
		Port:     util.Env("MAIL_SMTP_PORT", "465"),
	}

	smtpServer.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer.Host,
	}

	if err := smtpServer.setupClient(); err != nil {
		log.Fatalln("failed to initialize SMTP client: ", err)
		return nil
	}

	return smtpServer
}

// Send mail implements Sender.Send
func (smtpServer *SMTP) Send(content *Content) error {
	if err := smtpServer.send(content); err != nil {
		return nil
	}

	return nil
}

// setupClient is helper function to setup and create client
func (smtpServer *SMTP) setupClient() error {
	conn, err := tls.Dial("tcp", smtpServer.Host+":"+smtpServer.Port, smtpServer.TLSConfig)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, smtpServer.Host)
	if err != nil {
		return err
	}

	if err := smtpServer.withAuth(client); err != nil {
		return err
	}

	smtpServer.Client = client

	return nil
}

// withAuth will apply authentication for smtp client
func (smtpServer *SMTP) withAuth(client *smtp.Client) error {
	auth := smtp.PlainAuth("", smtpServer.From, smtpServer.Password, smtpServer.Host)

	if err := client.Auth(auth); err != nil {
		return err
	}

	if err := client.Mail(smtpServer.From); err != nil {
		return err
	}

	return nil
}

// send is helper function to send email
func (smtpServer *SMTP) send(content *Content) error {
	receivers := append(content.To, content.Cc...)
	receivers = append(receivers, content.Bcc...)

	for _, k := range receivers {
		if err := smtpServer.Client.Rcpt(k); err != nil {
			continue
		}
	}

	w, err := smtpServer.Client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(content.Construct())
	if err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return nil
	}

	smtpServer.Client.Quit()

	return nil
}
