package rmq

import (
	"github.com/semirm-dev/go-common/util"
	"github.com/streadway/amqp"
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
	ConsumerTag  string
	*Options
}

// Options struct
type Options struct {
	Queue     *QueueOpts
	Exchange  *ExchangeOpts
	QoS       *QoSOpts
	QueueBind *QueueBindOpts
	Consume   *ConsumeOpts
	Publish   *PublishOpts
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

// QoSOpts struct
type QoSOpts struct {
	PrefetchCount int
	PrefetchSize  int
	Global        bool
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
