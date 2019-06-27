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

## Usage

### Consumer

* **Create new rmq.Connection struct and assign valid callback for HandleMsgs**
```
consumer := &rmq.Connection{
    Config:     rmq.NewConfig(), // can be customized
    HandleMsgs: func(msgs <-chan amqp.Delivery) {
        for m := range msgs {
            log.Print(m)
        }
    },
    ResetSignal: make(chan int),
	EnabledHealthCheck: true, // false by default
}
```

* **Call setup func on created consumer (*this will setup conn, exchange, queue, etc...*)**
> Optionally, call ListenNotifyClose()
```
if err := consumer.Setup(); err != nil {
    return err
}

done := make(chan bool)

go consumer.ListenNotifyClose(done)
```

* **Call consume func on created consumer to start consuming messages**
```
go func() {
    if err := consumer.Consume(done); err != nil {
        log.Print("rmq consume error: ", err)
    }
}()

// NOTE: close done channel to close RMQ connection
// close(done)
<-done
```

* **Optionally, listen for reset signal from rmq connection and restart consumer.Consume()**
> NOTE: Move this part to core implementation
```
go func(done chan bool) {
	for {
		select {
		case s := <-consumer.ResetSignal:
			log.Print("consumer received rmq connection reset signal: ", s)

			// NOTE: required if we close done channel -> close(done)
			// done := make(chan bool)

			go func() {
				if err := consumer.Consume(done); err != nil {
					log.Print("rmq failed to consume: ", err)
					return
				}
			}()
		}
	}
}(done)
```

### Publisher

* **Create new rmq.Connection struct**
```
publisher := &rmq.Connection{
    Config:     rmq.NewConfig(), // can be customized
	ResetSignal: make(chan int),
}
```

* **Call setup func on created publisher (*this will setup conn, exchange, queue, etc...*)**
> Optionally, call ListenNotifyClose()
```
if err := publisher.Setup(); err != nil {
    return err
}

done := make(chan bool)

go publisher.ListenNotifyClose(done)
```

* **Optionally, listen for reset signal from rmq connection and re-create rmq connection**
> NOTE: Move this part to core implementation
```
go func() {
	for {
		select {
		case s := <-publisher.ResetSignal:
			log.Print("publisher received rmq connection reset signal: ", s)

			if err := publisher.Setup(); err != nil {
				return err
			}
		}
	}
}()
```

* **Optionally, set rmq headers**
```
publisher.WithHeaders(map[string]interface{}{
    "header-1": "value-1",
    "header-2": "value-2,
})
```

* **Publish message like so**
```
if err := publisher.Publish([]byte("message")); err != nil {
    log.Print("rmq publish error: ", err)
}
```

### Config customization, for both Consumer and Publisher (*pay close attention to options/structs nestings*)
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
