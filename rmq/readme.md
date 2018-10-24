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

* **Create new rmq.Consumer struct like so**
```
consumer := &rmq.Consumer{
    Config:     rmq.NewConfig(), // can be customized
    HandleMsgs: func(msgs <-chan amqp.Delivery) {
        for m := range msgs {
            log.Print(m)
        }
    },
}
```

* **Call setup func on created consumer (*this will setup conn, exchange, queue, etc...*)**
```
if err := consumer.Setup(); err != nil {
    return err
}
```

* **Call consume func on created consumer to start consuming messages**
```
done := make(chan bool)

go func() {
    if err := consumer.Consume(); err != nil {
        log.Print("consuming error: ", err)
    }
}()

<-done
```

### Publisher

* **Create new rmq.Publisher struct like so**
```
publisher := &rmq.Publisher{
    Config:     rmq.NewConfig(), // can be customized
}
```

* **Call setup func on created publisher (*this will setup conn, exchange, queue, etc...*)**
```
if err := publisher.Setup(); err != nil {
    return err
}
```

* **Optionally set rmq headers**
```
publisher.WithHeaders(map[string]interface{}{
    "header-1": "value-1",
    "header-2": "value-2,
})
```

* **Publish message like so**
```
if err := publisher.Publish([]byte("message")); err != nil {
    log.Print("publish error: ", err)
}
```

### Config customization, for both Consumer and Publisher (*pay close attention to options/structs nestings*)
```
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
		ConsumerOpts: &ConsumerOpts{
			&QueueOpts{
				Durable:          true,
				DeleteWhenUnused: false,
				Exclusive:        false,
				NoWait:           false,
				Args:             nil,
			},
			&ExchangeOpts{
				Durable:    true,
				AutoDelete: false,
				Internal:   false,
				NoWait:     false,
				Args:       nil,
			},
			&QueueBindOpts{
				NoWait: false,
				Args:   nil,
			},
			&ChannelConsumeOpts{
				AutoAck:   true,
				Exclusive: false,
				NoLocal:   false,
				NoWait:    false,
				Args:      nil,
			},
		},
		PublisherOpts: &PublisherOpts{
			&QueueOpts{
				Durable:          true,
				DeleteWhenUnused: false,
				Exclusive:        false,
				NoWait:           false,
				Args:             nil,
			},
			&ExchangeOpts{
				Durable:    true,
				AutoDelete: false,
				Internal:   false,
				NoWait:     false,
				Args:       nil,
			},
			&QueueBindOpts{
				NoWait: false,
				Args:   nil,
			},
			&ChannelPublishOpts{
				Mandatory: false,
				Immediate: false,
			},
		},
	}
}

```