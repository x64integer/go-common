## SMTP

* **ENV variables**

| ENV                | Default value  |
|:-------------------|:--------------:|
| MAIL_FROM          |                |
| MAIL_FROM_PASSWORD |                |
| MAIL_SMTP_HOST     | smtp.gmail.com |
| MAIL_SMTP_PORT     | 465            |

* **Usage**
```
smtpClient := mail.DefaultSMTP()

content := &mail.Content{
    To:         []string{"mail_1@gmail.com"},
    Cc:         []string{"mail_2@gmail.com"},
    Bcc:        []string{"mail_3@gmail.com"},
    Subject:    "Test mail",
    Body:       []byte("some mail body"),
    Attachment: []byte("some attachments"),
}

if err := smtpClient.Send(content); err != nil {
    log.Println("failed to send email: ", err)
}
```
