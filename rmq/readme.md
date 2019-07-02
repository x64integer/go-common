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
// config
config := rmq.NewConfig()
config.Exchange = "test_exchange"
config.Queue = "test_queue"
config.RoutingKey = "test_queue"

// setup connection
consumer := &rmq.Connection{
	Config: config,
	HandleMsgs: func(msgs <-chan amqp.Delivery) {
		for m := range msgs {
			log.Print(string(m.Body))
		}
	},
	ResetSignal: make(chan int),
}

if err := consumer.Setup(); err != nil {
	log.Fatal(err)
}

// start consumer
done := make(chan bool)

// optionally ListenNotifyClose and HandleResetSignalConsumer
go consumer.ListenNotifyClose(done)

go consumer.HandleResetSignalConsumer(done)

go func() {
	if err := consumer.Consume(done); err != nil {
		log.Print("rmq consume error: ", err)
	}
}()

<-done
```


### Publisher

```
// config
config := rmq.NewConfig()
config.Exchange = "test_exchange"
config.Queue = "test_queue"
config.RoutingKey = "test_queue"

// setup connection
publisher := &rmq.Connection{
	Config:      config,
	ResetSignal: make(chan int),
}

if err := publisher.Setup(); err != nil {
	log.Fatal(err)
}

// optionally ListenNotifyClose and HandleResetSignalPublisher
done := make(chan bool)

go publisher.ListenNotifyClose(done)

go publisher.HandleResetSignalPublisher(done)

// optionally set headers and publish message
publisher.WithHeaders(map[string]interface{}{
	"header-1": "value-1",
	"header-2": "value-2",
})

if err := publisher.Publish([]byte("message")); err != nil {
	log.Print("rmq publish error: ", err)
}

<-done
```

### Config customization, for both Consumer and Publisher
```
// NewConfig will initialize RMQ default config values
func NewConfig() *Config {
	return &Config{
		Host:         util.Env("RMQ_HOST", "localhost"),
		Port:         util.Env("RMQ_PORT", "5672"),
		Username:     util.Env("RMQ_USERNAME", "guest"),
		Password:     util.Env("RMQ_PASSWORD", "guest"),
		Exchange:     util.Env("RMQ_EXCHANGE", ""),
		ExchangeKind: util.Env("RMQ_EXCHANGE_KIND", "direct"),
		Queue:        util.Env("RMQ_QUEUE", ""),
		RoutingKey:   util.Env("RMQ_ROUTING_KEY", ""),
		ConsumerTag:  util.Env("RMQ_CONSUMER_TAG", ""),
		Options: &Options{
			Queue: &QueueOpts{
				Durable:          true,
				DeleteWhenUnused: false,
				Exclusive:        false,
				NoWait:           false,
				Args:             nil,
			},
			Exchange: &ExchangeOpts{
				Durable:    true,
				AutoDelete: false,
				Internal:   false,
				NoWait:     false,
				Args:       nil,
			},
			QoS: &QoSOpts{
				PrefetchCount: 1,
				PrefetchSize:  0,
				Global:        false,
			},
			QueueBind: &QueueBindOpts{
				NoWait: false,
				Args:   nil,
			},
			Consume: &ConsumeOpts{
				AutoAck:   true,
				Exclusive: false,
				NoLocal:   false,
				NoWait:    false,
				Args:      nil,
			},
			Publish: &PublishOpts{
				Mandatory: false,
				Immediate: false,
			},
		},
	}
}
```
