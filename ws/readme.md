## Usage

### Client

* **Create ws.Client**
```
client := &ws.Client{
	MessageHandler: &mHandler{},
}

done := make(chan bool)
ready := make(chan bool)

go client.Connect(done, ready, "ws://localhost:8080")

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
server := &ws.Server{
	Host:     "localhost",
	Port:     "8080",
	Endpoint: "/",
	MessageHandler: &mHandler{},
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
}
```
