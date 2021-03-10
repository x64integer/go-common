## ENV variables

| ENV                | Default value |
|:-------------------|:-------------:|
| RMQ_HOST           | localhost     |
| RMQ_PORT           | 5672          |
| RMQ_USERNAME       | guest         |
| RMQ_PASSWORD       | guest         |
| RMQ_EXCHANGE       |               |
| RMQ_EXCHANGE_KIND  | direct        |
| RMQ_QUEUE          |               |
| RMQ_ROUTING_KEY    |               |
| RMQ_CONSUMER_TAG   |               |

> TODO: Merge Consume and ConsumeWithConfig into single func

## Usage

### Consumer

```
credentials := rmq.NewCredentials()

// config
config := rmq.NewConfig()
config.Exchange = "test_exchange"
config.Queue = "test_queue"
config.RoutingKey = "test_queue"

// setup connection
consumer := &rmq.Connection{
    Credentials: credentials,
    Config: config,
    HandleMsg: func(msg <-chan amqp.Delivery) {
        for m := range msg {
            logrus.Info(config.Queue + " - " + string(m.Body))
        }
    },
    ResetSignal: make(chan int),
}

if err := consumer.Connect(true); err != nil {
    logrus.Fatal(err)
}

// start consumer
done := make(chan bool)

// optionally ListenNotifyClose and HandleResetSignalConsumer
go consumer.ListenNotifyClose(done)

go consumer.HandleResetSignalConsumer(done)

go func() {
    if err := consumer.Consume(done); err != nil {
        logrus.Error(err)
    }
}()

<-done
```


### Publisher

```
credentials := rmq.NewCredentials()

// config
config := rmq.NewConfig()
config.Exchange = "test_exchange"
config.Queue = "test_queue"
config.RoutingKey = "test_queue"

// setup connection
publisher := &rmq.Connection{
    Credentials: credentials,
    Config:      config,
    ResetSignal: make(chan int),
}

// pass true if there is only one publisher config
// else manually call publisher.ApplyConfig(*Config) for each configuration and
// call publisher.PublishWithConfig(*Config) if publisher.Config was not set!
if err := publisher.Connect(true); err != nil {
    logrus.Fatal(err)
}

// optionally ListenNotifyClose and HandleResetSignalPublisher
done := make(chan bool)

go publisher.ListenNotifyClose(done)

go publisher.HandleResetSignalPublisher(done)

wg := sync.WaitGroup{}

configB := rmq.NewConfig()
configB.Exchange = "test_exchange_b"
configB.Queue = "test_queue_b"
configB.RoutingKey = "test_queue_b"

if err := publisher.ApplyConfig(configB); err != nil {
    logrus.Error(err)
    return
}

for i := 0; i < 30000; i++ {
    go func() {
        wg.Add(1)
        defer wg.Done()

        if err := publisher.Publish([]byte(str.UUID())); err != nil {
            logrus.Error(err)
        }
    }()

    go func() {
        wg.Add(1)
        defer wg.Done()

        if err := publisher.PublishWithConfig(configB, []byte(str.UUID())); err != nil {
            logrus.Error(err)
        }
    }()
}

wg.Wait()

close(done)

<-done
```

> Config: https://github.com/semirm-dev/godev/blob/master/rmq/config.go
