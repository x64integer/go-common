## ENV variables

| ENV                          | Default value         |
|:-----------------------------|:---------------------:|
| MAIL_SERVICE_SMTP_HOST       | smtp.gmail.com        |
| MAIL_SERVICE_SMTP_PORT       | 465                   |
| MAIL_SERVICE_FROM            |                       |
| MAIL_SERVICE_FROM_PASSWORD   |                       |

* **Setup mail client**
```
smtpClient := mail.DefaultSMTP()

client := &mail.Client{
    Sender: smtpClient,
}
```

* **Construct mail content to be sent**
```
content := &mail.Content{
    To:         []string{"mail_1@gmail.com"},
    Cc:         []string{"mail_2@gmail.com"},
    Bcc:        []string{"mail_3@gmail.com"},
    Subject:    "Test mail",
    Content:    []byte("some content"),
    Attachment: []byte("some attachments"),
}
```

* **Send mail**
```
if err := client.Send(content); err != nil {
    log.Println("failed to send email: ", err)
}
```