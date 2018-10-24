package rmq

import (
	"github.com/streadway/amqp"
	"github.com/x64puzzle/go-common/util"
)

// Config for RMQ
type Config struct {
	Host         string
	Port         string
	Username     string
	Password     string
	Exchange     string
	ExchangeKind string
	Queue        string
	RoutingKey   string
	*ConsumerOpts
	*PublisherOpts
}

// ConsumerOpts struct
type ConsumerOpts struct {
	ConsumerTag string
	*QueueOpts
	*ExchangeOpts
	*QueueBindOpts
	*ChannelConsumeOpts
}

// PublisherOpts struct
type PublisherOpts struct {
	*QueueOpts
	*ExchangeOpts
	*QueueBindOpts
	*ChannelPublishOpts
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

// ExchangeOpts struct
type ExchangeOpts struct {
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

// QueueBindOpts struct
type QueueBindOpts struct {
	NoWait bool
	Args   amqp.Table
}

// ChannelConsumeOpts struct
type ChannelConsumeOpts struct {
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}

// ChannelPublishOpts struct
type ChannelPublishOpts struct {
	Mandatory bool
	Immediate bool
}

// NewConfig will initialize RMQ config for error publisher
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
		ConsumerOpts: &ConsumerOpts{
			ConsumerTag: util.Env("RMQ_CONSUMER_TAG", ""),
			QueueOpts: &QueueOpts{
				Durable:          true,
				DeleteWhenUnused: false,
				Exclusive:        false,
				NoWait:           false,
				Args:             nil,
			},
			ExchangeOpts: &ExchangeOpts{
				Durable:    true,
				AutoDelete: false,
				Internal:   false,
				NoWait:     false,
				Args:       nil,
			},
			QueueBindOpts: &QueueBindOpts{
				NoWait: false,
				Args:   nil,
			},
			ChannelConsumeOpts: &ChannelConsumeOpts{
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
