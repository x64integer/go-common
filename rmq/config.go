package rmq

import (
	"github.com/semirm-dev/go-dev/env"
	"github.com/streadway/amqp"
)

// Config for RMQ
type Config struct {
	Host     string
	Port     string
	Username string
	Password string

	Exchange     string
	ExchangeKind string

	Queue       string
	RoutingKey  string
	ConsumerTag string

	*Options
}

// Options struct
type Options struct {
	Exchange *ExchangeOpts
	QoS      *QoSOpts

	Queue     *QueueOpts
	QueueBind *QueueBindOpts

	Consume *ConsumeOpts
	Publish *PublishOpts
}

// ExchangeOpts struct
type ExchangeOpts struct {
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

// QoSOpts struct
type QoSOpts struct {
	PrefetchCount int
	PrefetchSize  int
	Global        bool
}

// QueueOpts struct
type QueueOpts struct {
	Durable          bool
	DeleteWhenUnused bool
	Exclusive        bool
	Internal         bool
	NoWait           bool
	Args             amqp.Table
}

// QueueBindOpts struct
type QueueBindOpts struct {
	NoWait bool
	Args   amqp.Table
}

// ConsumeOpts struct
type ConsumeOpts struct {
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}

// PublishOpts struct
type PublishOpts struct {
	Mandatory bool
	Immediate bool
}

// NewConfig will initialize RMQ default config values
func NewConfig() *Config {
	return &Config{
		Host:     env.Get("RMQ_HOST", "localhost"),
		Port:     env.Get("RMQ_PORT", "5672"),
		Username: env.Get("RMQ_USERNAME", "guest"),
		Password: env.Get("RMQ_PASSWORD", "guest"),

		Exchange:     env.Get("RMQ_EXCHANGE", ""),
		ExchangeKind: env.Get("RMQ_EXCHANGE_KIND", "direct"),

		Queue:       env.Get("RMQ_QUEUE", ""),
		RoutingKey:  env.Get("RMQ_ROUTING_KEY", ""),
		ConsumerTag: env.Get("RMQ_CONSUMER_TAG", ""),

		Options: &Options{
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

			Queue: &QueueOpts{
				Durable:          true,
				DeleteWhenUnused: false,
				Exclusive:        false,
				NoWait:           false,
				Args:             nil,
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
