package rmq

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/streadway/amqp"
)

const (
	Reconnected = iota
)

var (
	// reconnectTime is default time to wait for rmq reconnect on Conn.NotifyClose() event - situation when rmq sends signal about shutdown
	reconnectTime = 20 * time.Second
)

// Connection for RMQ
type Connection struct {
	Credentials *Credentials
	Config      *Config

	// amqp
	Conn        *amqp.Connection
	Channel     *amqp.Channel
	Headers     amqp.Table
	ContentType string

	// connection reset
	ResetSignal   chan int
	ReconnectTime time.Duration
	Retrying      bool

	// callbacks
	HandleMsg                  func(msg <-chan amqp.Delivery)
	HandleResetSignalConsumer  func(chan bool)
	HandleResetSignalPublisher func(chan bool)
}

// Connect to RabbitMQ and initialize channel
func (c *Connection) Connect(applyConfig bool) error {
	if c.Credentials == nil {
		return errors.New("invalid/nil Credentials")
	}

	c.applyDefaults()

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", c.Credentials.Username, c.Credentials.Password, c.Credentials.Host, c.Credentials.Port))
	if err != nil {
		return err
	}
	c.Conn = conn

	if applyConfig {
		if err := c.ApplyConfig(c.Config); err != nil {
			return err
		}
	}

	return nil
}

// ApplyConfig will initialize channel, exchange, qos and bind queues
// RabbitMQ declarations
func (c *Connection) ApplyConfig(config *Config) error {
	if c.Conn == nil {
		return errors.New("amqp connection not initialized")
	}

	if config == nil {
		return errors.New("invalid/nil Config")
	}

	ch, err := c.Conn.Channel()
	if err != nil {
		return err
	}

	c.Channel = ch

	if err := c.exchangeDeclare(config.Exchange, config.ExchangeKind, config.Options.Exchange); err != nil {
		return err
	}

	if err := c.qos(config.Options.QoS); err != nil {
		return err
	}

	if _, err := c.queueDeclare(config.Queue, config.Options.Queue); err != nil {
		return err
	}

	if err := c.queueBind(config.Queue, config.RoutingKey, config.Exchange, config.Options.QueueBind); err != nil {
		return err
	}

	return nil
}

// Consume data from RMQ
func (c *Connection) Consume(done chan bool) error {
	msg, err := c.Channel.Consume(
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

	go c.HandleMsg(msg)

	logrus.Info("waiting for messages...")

	for {
		select {
		case <-done:
			if err := c.Channel.Close(); err != nil {
				logrus.Error("failed to close channel: ", err.Error())
				return err
			}

			if err := c.Conn.Close(); err != nil {
				logrus.Error("failed to close connection: ", err.Error())
				return err
			}

			return nil
		}
	}
}

// Publish payload to RMQ
func (c *Connection) Publish(payload []byte) error {
	if c.Config == nil {
		return errors.New("invalid/nil Config")
	}

	err := c.Channel.Publish(
		c.Config.Exchange,
		c.Config.RoutingKey,
		c.Config.Options.Publish.Mandatory,
		c.Config.Options.Publish.Immediate,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  c.ContentType,
			Body:         payload,
			Headers:      c.Headers,
		})

	return err
}

func (c *Connection) PublishWithConfig(config *Config, payload []byte) error {
	err := c.Channel.Publish(
		config.Exchange,
		config.RoutingKey,
		config.Options.Publish.Mandatory,
		config.Options.Publish.Immediate,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  c.ContentType,
			Body:         payload,
			Headers:      c.Headers,
		})

	return err
}

// WithHeaders will set headers to be sent
func (c *Connection) WithHeaders(h amqp.Table) *Connection {
	c.Headers = h

	return c
}

// ListenNotifyClose will listen for rmq connection shutdown and attempt to re-create rmq connection
func (c *Connection) ListenNotifyClose(done chan bool) {
	connClose := make(chan *amqp.Error)
	c.Conn.NotifyClose(connClose)

	go func() {
		for {
			select {
			case err := <-connClose:
				logrus.Warn("rmq connection lost: ", err)
				logrus.Warn("reconnecting to rmq in ", c.ReconnectTime.String())

				c.Retrying = true

				time.Sleep(c.ReconnectTime)

				if err := c.recreateConn(); err != nil {
					killService("failed to recreate rmq connection: ", err)
				}

				logrus.Infof("sending signal %v to rmq connection", Reconnected)

				c.ResetSignal <- Reconnected

				logrus.Infof("signal %v sent to rmq connection", Reconnected)

				// important step!
				// recreate connClose channel so we can listen for NotifyClose once again
				connClose = make(chan *amqp.Error)
				c.Conn.NotifyClose(connClose)

				c.Retrying = false
			}
		}
	}()

	<-done
}

// queueDeclare is helper function to declare queue
func (c *Connection) queueDeclare(name string, opts *QueueOpts) (amqp.Queue, error) {
	queue, err := c.Channel.QueueDeclare(
		name,
		opts.Durable,
		opts.DeleteWhenUnused,
		opts.Exclusive,
		opts.NoWait,
		opts.Args,
	)

	return queue, err
}

// exchangeDeclare is helper function to declare exchange
func (c *Connection) exchangeDeclare(name string, kind string, opts *ExchangeOpts) error {
	err := c.Channel.ExchangeDeclare(
		name,
		kind,
		opts.Durable,
		opts.AutoDelete,
		opts.Internal,
		opts.NoWait,
		opts.Args,
	)

	return err
}

// qos is helper function to define QoS for channel
func (c *Connection) qos(opts *QoSOpts) error {
	err := c.Channel.Qos(
		opts.PrefetchCount,
		opts.PrefetchSize,
		opts.Global,
	)

	return err
}

// queueBind is helper function to bind queue to exchange
func (c *Connection) queueBind(queue string, routingKey string, exchange string, opts *QueueBindOpts) error {
	err := c.Channel.QueueBind(
		queue,
		routingKey,
		exchange,
		opts.NoWait,
		opts.Args,
	)

	return err
}

// applyDefaults is helper function to setup some default Connection properties
func (c *Connection) applyDefaults() {
	if c.ReconnectTime == 0 {
		c.ReconnectTime = reconnectTime
	}

	if c.HandleResetSignalConsumer == nil {
		c.HandleResetSignalConsumer = c.handleResetSignalConsumer
	}

	if c.HandleResetSignalPublisher == nil {
		c.HandleResetSignalPublisher = c.handleResetSignalPublisher
	}

	if c.ContentType == "" {
		c.ContentType = "text/plain"
	}
}

// handleResetSignalConsumer is default callback for consumer to run when reset signal occurs
func (c *Connection) handleResetSignalConsumer(done chan bool) {
	go func(done chan bool) {
		for {
			select {
			case s := <-c.ResetSignal:
				logrus.Warn("consumer received rmq connection reset signal: ", s)

				if done == nil {
					done = make(chan bool)
				}

				go func() {
					if err := c.Consume(done); err != nil {
						logrus.Fatal("rmq failed to consume: ", err)
					}
				}()
			}
		}
	}(done)

	<-done
}

// handleResetSignalPublisher is default callback for publisher to run when reset signal occurs
func (c *Connection) handleResetSignalPublisher(done chan bool) {
	go func() {
		for {
			select {
			case s := <-c.ResetSignal:
				logrus.Warn("publisher received rmq connection reset signal: ", s)
			}
		}
	}()

	<-done
}

// recreateConn for rmq
func (c *Connection) recreateConn() error {
	logrus.Info("trying to recreate rmq connection for host: ", c.Credentials.Host)

	return c.Connect(true)
}

// killService with message passed to console output
func killService(msg ...interface{}) {
	logrus.Warn(msg...)
	os.Exit(101)
}
