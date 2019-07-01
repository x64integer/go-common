## ENV variables

| ENV    | Default value |
|:-------|:-------------:|
| WS_URL |               |

## Usage

### Client

* **Create ws.Client**
```
config := ws.NewConfig()
config.WSURL = "ws://localhost:8080"

client := &ws.Client{
	Config:         config,
	MessageHandler: &mHandler{},
}

done := make(chan bool)
ready := make(chan bool)

go client.Connect(done, ready)

<-ready

go func() {
	for i := 0; i < 50000; i++ {
		// ready to send messages to websocket channel
		if err := client.SendText([]byte(fmt.Sprint("message: ", i))); err != nil {
			logrus.Fatal("failed to send message: ", err)
		}
	}
}()

// go func() {
// 	select {
// 	case <-time.After(2 * time.Second):
// 		close(done)
// 	}
// }()

<-done
```

### Server

* **Create ws.Server**
```
config := ws.NewConfig()

h := &mHandler{}

server := &ws.Server{
	Config:         config,
	MessageHandler: h,
}

done := make(chan bool)

go server.Run(done)

<-done
```

### Handler example
```
// mHandler example impl
type mHandler struct{}

func (h *mHandler) OnMessage(in []byte, reply func(int, []byte) error) {
	// handle message from ws channel
	logrus.Info("received: " + string(in))

	reply(ws.TextMessage, []byte("reply from server: "+string(in)))
}

func (h *mHandler) OnError(err error) {
	// handle error from ws channel
	logrus.Error("error from ws connection: ", err)

	if connClosed, ok := err.(*ws.ConnectionClosed); ok {
		logrus.Error(connClosed)
	}

	if strings.Contains(err.Error(), "closed network connection") {
		logrus.Error(err.Error())
	}
}
```
