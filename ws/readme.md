## ENV variables

| ENV    | Default value |
|:-------|:-------------:|
| WS_URL |               |

## Usage

* **Create ws.Connection**
```
c := &ws.Connection{
    Config: ws.NewConfig(), // can be customized
    OnMessage: func(in []byte) {
        // handle message from ws channel
    },
    OnError: func(err error) {
        // handle error from ws channel
    },
}
```

* **Setup ws connection and listen for messages**
```
if err := c.Setup(); err != nil {
    log.Print("error on ws connection setup: ", err)
}
```

* **Send message to ws channel**
```
// text type
c.SendText([]byte("message"))
// binary type
c.SendBinary([]byte("message"))
```