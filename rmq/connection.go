package rmq

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/streadway/amqp"
)

var (
	// reconnectTime is default time to wait for rmq reconnect on Conn.NotifyClose() event - situation when rmq sends signal about shutdown
	reconnectTime = 20 * time.Second
)

// Connection for RMQ
type Connection struct {
	Config *Config

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
	HandleMsgs                 func(msgs <-chan amqp.Delivery)
	HandleResetSignalConsumer  func(chan bool)
	HandleResetSignalPublisher func(chan bool)
}

// Connect to RabbitMQ and initialize channel
func (c *Connection) Connect(declareChannel bool) error {
	if c.Config == nil {
		return errors.New("nil Config struct for RMQ Connection -> make sure valid Config is accessible to Connection")
	}

	c.applyDefaults()

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", c.Config.Username, c.Config.Password, c.Config.Host, c.Config.Port))
	if err != nil {
		return err
	}
	c.Conn = conn

	if declareChannel {
		if err := c.declareChannel(); err != nil {
			return err
		}
	}

	return nil
}

// declareChannel will initialize channel, exchange, qos and bind queues
// RabbitMQ declarations
func (c *Connection) declareChannel() error {
	if c.Conn == nil {
		return errors.New("amqp connection not initialized")
	}

	ch, err := c.Conn.Channel()
	if err != nil {
		return err
	}

	c.Channel = ch

	if err := c.exchangeDeclare(c.Config.Exchange, c.Config.ExchangeKind, c.Config.Options.Exchange); err != nil {
		return err
	}

	if err := c.qos(c.Config.Options.QoS); err != nil {
		return err
	}

	if _, err := c.queueDeclare(c.Config.Queue, c.Config.Options.Queue); err != nil {
		return err
	}

	if err := c.queueBind(c.Config.Queue, c.Config.RoutingKey, c.Config.Exchange, c.Config.Options.QueueBind); err != nil {
		return err
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

	logrus.Info("waiting for messages...")

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

// PublishWithKey will publish payload to RMQ using given routingKey instead of routingKey provided in *Config
func (c *Connection) PublishWithKey(routingKey string, payload []byte) error {
	err := c.Channel.Publish(
		c.Config.Exchange,
		routingKey,
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

// PublishWithExchange will publish payload to RMQ using given exchange and routingKey instead of exchange/routingKey provided in *Config
func (c *Connection) PublishWithExchange(exchange, routingKey string, payload []byte) error {
	err := c.Channel.Publish(
		exchange,
		routingKey,
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

// DeclareWithConfig will initialize additional queues and exchanges on existing rmq setup/channel
func (c *Connection) DeclareWithConfig(config []*Config) error {
	if c.Channel == nil {
		ch, err := c.Conn.Channel()
		if err != nil {
			return err
		}

		c.Channel = ch
	}

	for _, conf := range config {
		if err := c.exchangeDeclare(conf.Exchange, conf.ExchangeKind, conf.Options.Exchange); err != nil {
			return err
		}

		if err := c.qos(conf.Options.QoS); err != nil {
			return err
		}

		if _, err := c.queueDeclare(conf.Queue, conf.Options.Queue); err != nil {
			return err
		}

		if err := c.queueBind(conf.Queue, conf.RoutingKey, conf.Exchange, conf.Options.QueueBind); err != nil {
			return err
		}
	}

	return nil
}

// ConsumeWithConfig will start consumer with passed config values
func (c *Connection) ConsumeWithConfig(done chan bool, config *Config, callback func(msgs <-chan amqp.Delivery)) error {
	msgs, err := c.Channel.Consume(
		config.Queue,
		config.ConsumerTag,
		config.Options.Consume.AutoAck,
		config.Options.Consume.Exclusive,
		config.Options.Consume.NoLocal,
		config.Options.Consume.NoWait,
		config.Options.Consume.Args,
	)
	if err != nil {
		return err
	}

	go callback(msgs)

	logrus.Info("waiting for messages...")

	for {
		select {
		case <-done:
			c.Channel.Close()
			c.Conn.Close()

			return nil
		}
	}
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

				logrus.Info("sending signal 1 to rmq connection")

				c.ResetSignal <- 1

				logrus.Info("signal 1 sent to rmq connection")

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
	logrus.Info("trying to recreate rmq connection for host: ", c.Config.Host)

	return c.Connect(true)
}

// killService with message passed to console output
func killService(msg ...interface{}) {
	logrus.Warn(msg...)
	os.Exit(101)
}
