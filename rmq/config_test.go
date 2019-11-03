package rmq_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/semirm-dev/go-dev/env"
	"github.com/semirm-dev/go-dev/rmq"
)

func TestNewConfig(t *testing.T) {
	expected := &rmq.Config{
		Host:         env.Get("RMQ_HOST", "localhost"),
		Port:         env.Get("RMQ_PORT", "5672"),
		Username:     env.Get("RMQ_USERNAME", "guest"),
		Password:     env.Get("RMQ_PASSWORD", "guest"),
		Exchange:     env.Get("RMQ_EXCHANGE", ""),
		ExchangeKind: env.Get("RMQ_EXCHANGE_KIND", "direct"),
		Queue:        env.Get("RMQ_QUEUE", ""),
		RoutingKey:   env.Get("RMQ_ROUTING_KEY", ""),
		ConsumerTag:  env.Get("RMQ_CONSUMER_TAG", ""),
		Options: &rmq.Options{
			Queue: &rmq.QueueOpts{
				Durable:          true,
				DeleteWhenUnused: false,
				Exclusive:        false,
				NoWait:           false,
				Args:             nil,
			},
			Exchange: &rmq.ExchangeOpts{
				Durable:    true,
				AutoDelete: false,
				Internal:   false,
				NoWait:     false,
				Args:       nil,
			},
			QoS: &rmq.QoSOpts{
				PrefetchCount: 1,
				PrefetchSize:  0,
				Global:        false,
			},
			QueueBind: &rmq.QueueBindOpts{
				NoWait: false,
				Args:   nil,
			},
			Consume: &rmq.ConsumeOpts{
				AutoAck:   true,
				Exclusive: false,
				NoLocal:   false,
				NoWait:    false,
				Args:      nil,
			},
			Publish: &rmq.PublishOpts{
				Mandatory: false,
				Immediate: false,
			},
		},
	}

	config := rmq.NewConfig()

	assert.Equal(t, expected, config)
}
