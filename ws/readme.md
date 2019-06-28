## ENV variables

| ENV    | Default value |
|:-------|:-------------:|
| WS_URL |               |

## Usage

### Client

* **Create ws.Client**
```
config := ws.NewConfig()
config.WSURL = "ws://echo.websocket.org"

client := &ws.Client{
	Config:       config,
	EventHandler: &mHandler{},
}

done := make(chan bool)
ready := make(chan bool)

go client.Run(done, ready)

<-ready

// ready to send messages to websocket channel
if err := client.SendText([]byte("test message")); err != nil {
	logrus.Fatal("failed to send message: ", err)
}
```

### Server (in progress)
> NOTE: need to improve OnMessage handler (there is no reply to client support at the moment)

* **Create ws.Server**
```
config := ws.NewConfig()
config.Endpoint = "/test"

server := &ws.Server{
	Config:       config,
	EventHandler: &mHandler{},
}

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