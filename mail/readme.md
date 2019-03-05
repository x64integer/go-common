## ENV variables

| ENV                          | Default value         |
|:-----------------------------|:---------------------:|
| MAIL_SERVICE_SMTP_HOST       | smtp.gmail.com        |
| MAIL_SERVICE_SMTP_PORT       | 465                   |
| MAIL_SERVICE_FROM            |                       |
| MAIL_SERVICE_FROM_PASSWORD   |                       |

* **Create mail entity to send**
```
m := &mail.Entity{
    To:         []string{"mail_1@mail.com"},
    Cc:         []string{"mail_2@mail.com"},
    Bcc:        []string{"mail_3@mail.com"},
    Content:    []byte("some content"),
    Attachment: []byte("some attachments"),
}
```

* **Send mail**
```
if err := m.Send(); err != nil {
    log.Println("failed to send mail: ", err)
}
```

### TODO
* Abstract mail interface
* Make it possible to use different mail services (custom smtp, sendgrid...) at any time