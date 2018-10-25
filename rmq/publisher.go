package rmq

import (
	"errors"
	"fmt"

	"github.com/streadway/amqp"
)

// Publisher for RMQ
type Publisher struct {
	Config  *Config
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Headers amqp.Table
}

// Setup RMQ Publisher
func (p *Publisher) Setup() error {
	if p.Config == nil {
		return errors.New("nil Config struct for RMQ Publisher -> make sure valid Config is accessible to Publisher")
	}

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", p.Config.Username, p.Config.Password, p.Config.Host, p.Config.Port))
	if err != nil {
		return err
	}
	p.Conn = conn

	ch, err := p.Conn.Channel()
	if err != nil {
		return err
	}

	p.Channel = ch

	_, err = p.Channel.QueueDeclare(
		p.Config.Queue,
		p.Config.PublisherOpts.QueueOpts.Durable,
		p.Config.PublisherOpts.QueueOpts.DeleteWhenUnused,
		p.Config.PublisherOpts.QueueOpts.Exclusive,
		p.Config.PublisherOpts.QueueOpts.NoWait,
		p.Config.PublisherOpts.QueueOpts.Args,
	)
	if err != nil {
		return err
	}

	if err := p.Channel.ExchangeDeclare(
		p.Config.Exchange,
		p.Config.ExchangeKind,
		p.Config.PublisherOpts.ExchangeOpts.Durable,
		p.Config.PublisherOpts.ExchangeOpts.AutoDelete,
		p.Config.PublisherOpts.ExchangeOpts.Internal,
		p.Config.PublisherOpts.ExchangeOpts.NoWait,
		p.Config.PublisherOpts.ExchangeOpts.Args,
	); err != nil {
		return err
	}

	if err := p.Channel.QueueBind(
		p.Config.Queue,
		p.Config.RoutingKey,
		p.Config.Exchange,
		p.Config.PublisherOpts.QueueBindOpts.NoWait,
		p.Config.PublisherOpts.QueueBindOpts.Args,
	); err != nil {
		return err
	}

	return nil
}

// Publish payload to RMQ
func (p *Publisher) Publish(payload []byte) error {
	if err := p.Channel.Publish(
		p.Config.Exchange,
		p.Config.RoutingKey,
		p.Config.PublisherOpts.PublishOpts.Mandatory,
		p.Config.PublisherOpts.PublishOpts.Immediate,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         payload,
			Headers:      p.Headers,
		}); err != nil {
		return err
	}

	defer p.Conn.Close()
	defer p.Channel.Close()

	return nil
}

// WithHeaders will set headers to be sent
func (p *Publisher) WithHeaders(h amqp.Table) *Publisher {
	p.Headers = h

	return p
}
