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
        log.Println("error from ws connection: ", err)

        if strings.Contains(err.Error(), "closed network connection") {
			os.Exit(1)
		}
    },
    OnConnClose: func(code int, msg string) {
        // handle closed ws connection
        log.Printf("closed conn: code=%v - msg=%v\n", code, msg)

        os.Exit(1)
    },
}
```

* **Setup ws connection and listen for messages**
```
if err := c.Setup(); err != nil {
    log.Print("ws connection setup error: ", err)
}
```

* **Send message to ws channel**
```
// text type
c.SendText([]byte("message"))
// binary type
c.SendBinary([]byte("message"))
```