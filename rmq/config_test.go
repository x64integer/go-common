package rmq_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/semirm-dev/go-common/rmq"
	"github.com/semirm-dev/go-common/util"
)

func TestNewConfig(t *testing.T) {
	expected := &rmq.Config{
		Host:         util.Env("RMQ_HOST", "localhost"),
		Port:         util.Env("RMQ_PORT", "5672"),
		Username:     util.Env("RMQ_USERNAME", "guest"),
		Password:     util.Env("RMQ_PASSWORD", "guest"),
		Exchange:     util.Env("RMQ_EXCHANGE", ""),
		ExchangeKind: util.Env("RMQ_EXCHANGE_KIND", "direct"),
		Queue:        util.Env("RMQ_QUEUE", ""),
		RoutingKey:   util.Env("RMQ_ROUTING_KEY", ""),
		ConsumerTag:  util.Env("RMQ_CONSUMER_TAG", ""),
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
