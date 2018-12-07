package rmq

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

// reconnectTime is default time to wait for rmq reconnect on Conn.NotifyClose() event - situation when rmq sends signal about shutdown
var reconnectTime = 15 * time.Second

// Connection for RMQ
type Connection struct {
	Config        *Config
	Conn          *amqp.Connection
	Channel       *amqp.Channel
	HandleMsgs    func(msgs <-chan amqp.Delivery)
	Headers       amqp.Table
	ResetSignal   chan int
	ReconnectTime time.Duration
}

// Setup RMQ Connection
func (c *Connection) Setup() error {
	if c.Config == nil {
		return errors.New("nil Config struct for RMQ Connection -> make sure valid Config is accessible to Connection")
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
		c.Config.Options.Queue.Durable,
		c.Config.Options.Queue.DeleteWhenUnused,
		c.Config.Options.Queue.Exclusive,
		c.Config.Options.Queue.NoWait,
		c.Config.Options.Queue.Args,
	)
	if err != nil {
		return err
	}

	if err := c.Channel.ExchangeDeclare(
		c.Config.Exchange,
		c.Config.ExchangeKind,
		c.Config.Options.Exchange.Durable,
		c.Config.Options.Exchange.AutoDelete,
		c.Config.Options.Exchange.Internal,
		c.Config.Options.Exchange.NoWait,
		c.Config.Options.Exchange.Args,
	); err != nil {
		return err
	}

	if err := c.Channel.QueueBind(
		c.Config.Queue,
		c.Config.RoutingKey,
		c.Config.Exchange,
		c.Config.Options.QueueBind.NoWait,
		c.Config.Options.QueueBind.Args,
	); err != nil {
		return err
	}

	if c.ReconnectTime == 0 {
		c.ReconnectTime = reconnectTime
	}

	return nil
}

// Consume data from RMQ
func (c *Connection) Consume(done chan bool) error {
	msgs, err := c.Channel.Consume(
		c.Config.Queue,
		c.Config.ConsumerTag,
		c.Config.Options.Consume.AutoAck,
		c.Config.Options.Consume.Exclusive,
		c.Config.Options.Consume.NoLocal,
		c.Config.Options.Consume.NoWait,
		c.Config.Options.Consume.Args,
	)
	if err != nil {
		return err
	}

	go c.HandleMsgs(msgs)

	log.Print("Waiting for messages...")

	for {
		select {
		case <-done:
			c.Channel.Close()
			c.Conn.Close()

			return nil
		}
	}
}

// Publish payload to RMQ
func (c *Connection) Publish(payload []byte) error {
	if err := c.Channel.Publish(
		c.Config.Exchange,
		c.Config.RoutingKey,
		c.Config.Options.Publish.Mandatory,
		c.Config.Options.Publish.Immediate,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         payload,
			Headers:      c.Headers,
		}); err != nil {
		return err
	}

	return nil
}

// WithHeaders will set headers to be sent
func (c *Connection) WithHeaders(h amqp.Table) *Connection {
	c.Headers = h

	return c
}

// ListenNotifyClose will listen for rmq connection shutdown and attempt to re-create rmq connection
func (c *Connection) ListenNotifyClose() {
	connClose := make(chan *amqp.Error)
	c.Conn.NotifyClose(connClose)

	go func() {
		err := <-connClose

		log.Print("rmq connection lost: ", err)
		log.Printf("reconnecting to rmq in %v...", c.ReconnectTime.String())

		select {
		case <-time.After(c.ReconnectTime):
			log.Print("re-creating rmq connection")

			if err := c.Setup(); err != nil {
				log.Print("failed to recreate rmq connection: ", err)

				os.Exit(101)
			}
		}

		c.ResetSignal <- 1

		log.Print("rmq reconnection successul, signal 1 sent")
	}()
}
