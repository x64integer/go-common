## ENV variables

| ENV    | Default value |
|:-------|:-------------:|
| WS_URL |               |

## Usage

### Client

* **Create ws.Client**
```
client := &ws.Client{
    Config:       ws.NewConfig(), // can be customized
    EventHandler: &mHandler{},
    // DisabledReader: true, // false by default, if set to true client will not read messages from websocket channel
}
```

* **Setup ws Client and listen for messages**
```
done := make(chan bool)
ready := make(chan bool)

go client.Run(done, ready)

<-ready

// ready to send messages to websocket channel
```

* **Send message to ws channel**
```
// text type
client.SendText([]byte("message"))
// binary type
client.SendBinary([]byte("message"))
```

### Server (in progress)
> NOTE: need to improve OnMessage handler (there is no reply to client support at the moment)

* **Create ws.Server**
```
config := ws.NewConfig()
// http server configurations
config.Endpoint = "/test"
config.Host = "localhost"
config.Port = "8080"

server := &ws.Server{
    Config:       config,
    EventHandler: &mHandler{},
    // DisabledReader: true, // false by default, if set to true client will not read messages from websocket channel
}
```

* **Setup ws server and listen for messages**
```
done := make(chan bool)

go server.Run(done)

<-done
```

* **Handler example**
```
// mHandler example impl
type mHandler struct{}

func (h *mHandler) OnMessage(in []byte) {
	// handle message from ws channel
	log.Println("received: " + string(in))
}

func (h *mHandler) OnError(err error) {
	// handle error from ws channel
	log.Println("error from ws connection: ", err)

	if connClosed, ok := err.(*ws.ConnectionClosed); ok {
		log.Println(connClosed)
	}

	if strings.Contains(err.Error(), "closed network connection") {
		os.Exit(1)
	}
}
```