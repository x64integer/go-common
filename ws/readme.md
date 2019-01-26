## ENV variables

| ENV    | Default value |
|:-------|:-------------:|
| WS_URL |               |

## Usage

### Client

* **Create ws.Client**
```
c := &ws.Client{
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

* **Setup ws Client and listen for messages**
```
done := make(chan bool)

go c.Setup(done)

<-done
```

* **Send message to ws channel**
```
// text type
c.SendText([]byte("message"))
// binary type
c.SendBinary([]byte("message"))
```

### Server (in progress)
> NOTE: need to improve OnMessage handler (there is no reply to client support at the moment)

* **Create ws.Server**
```
config := ws.NewConfig()
// http server configurations
config.Endpoint = "/test"
config.Host = "localhost"
config.Port = "8080

s := &ws.Server{
    Config: config,
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

* **Setup ws server and listen for messages**
```
done := make(chan bool)

go s.Run(done)

<-done
```