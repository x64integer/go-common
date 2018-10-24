package rmq

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

// Consumer for RMQ
type Consumer struct {
	Config     *Config
	Conn       *amqp.Connection
	Channel    *amqp.Channel
	HandleMsgs func(msgs <-chan amqp.Delivery)
}

// Setup RMQ Consumer
func (c *Consumer) Setup() error {
	if c.Config == nil {
		return errors.New("nil Config struct for RMQ Consumer -> make sure valid Config is accessible to Consumer")
	}

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", c.Config.Username, c.Config.Password, c.Config.Host, c.Config.Port))
	if err != nil {
		return err
	}
	c.Conn = conn

	ch, err := c.Conn.Channel()
	if err != nil {
		return err
	}

	c.Channel = ch

	_, err = c.Channel.QueueDeclare(
		c.Config.Queue,
		c.Config.ConsumerOpts.QueueOpts.Durable,
		c.Config.ConsumerOpts.QueueOpts.DeleteWhenUnused,
		c.Config.ConsumerOpts.QueueOpts.Exclusive,
		c.Config.ConsumerOpts.QueueOpts.NoWait,
		c.Config.ConsumerOpts.QueueOpts.Args,
	)
	if err != nil {
		return err
	}

	if err := c.Channel.ExchangeDeclare(
		c.Config.Exchange,
		c.Config.ExchangeKind,
		c.Config.ConsumerOpts.ExchangeOpts.Durable,
		c.Config.ConsumerOpts.ExchangeOpts.AutoDelete,
		c.Config.ConsumerOpts.ExchangeOpts.Internal,
		c.Config.ConsumerOpts.ExchangeOpts.NoWait,
		c.Config.ConsumerOpts.ExchangeOpts.Args,
	); err != nil {
		return err
	}

	if err := c.Channel.QueueBind(
		c.Config.Queue,
		c.Config.RoutingKey,
		c.Config.Exchange,
		c.Config.ConsumerOpts.QueueBindOpts.NoWait,
		c.Config.ConsumerOpts.QueueBindOpts.Args,
	); err != nil {
		return err
	}

	return nil
}

// Consume data from RMQ
func (c *Consumer) Consume() error {
	msgs, err := c.Channel.Consume(
		c.Config.Queue,
		c.Config.ConsumerOpts.ConsumerTag,
		c.Config.ConsumerOpts.ChannelConsumeOpts.AutoAck,
		c.Config.ConsumerOpts.ChannelConsumeOpts.Exclusive,
		c.Config.ConsumerOpts.ChannelConsumeOpts.NoLocal,
		c.Config.ConsumerOpts.ChannelConsumeOpts.NoWait,
		c.Config.ConsumerOpts.ChannelConsumeOpts.Args,
	)
	if err != nil {
		return err
	}

	defer c.Conn.Close()
	defer c.Channel.Close()

	killer := make(chan *amqp.Error)

	go func() {
		err := <-killer

		log.Print("consume connection error: ", err)

		os.Exit(101)
	}()

	c.Conn.NotifyClose(killer)

	done := make(chan bool)

	go c.HandleMsgs(msgs)

	log.Print("Waiting for messages...")

	<-done

	return nil
}
