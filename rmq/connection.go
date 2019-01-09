package rmq

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/streadway/amqp"
)

var (
	// reconnectTime is default time to wait for rmq reconnect on Conn.NotifyClose() event - situation when rmq sends signal about shutdown
	reconnectTime = 20 * time.Second
	// healthCheckTime is time interval for healthCheck
	healthCheckTime = 5 * time.Second
)

// Connection for RMQ
type Connection struct {
	Config              *Config
	Conn                *amqp.Connection
	Channel             *amqp.Channel
	HandleMsgs          func(msgs <-chan amqp.Delivery)
	Headers             amqp.Table
	ResetSignal         chan int
	ReconnectTime       time.Duration
	Retrying            bool
	DisabledHealthCheck bool
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

	if err := c.Channel.Qos(
		c.Config.Options.QoS.PrefetchCount,
		c.Config.Options.QoS.PrefetchSize,
		c.Config.Options.QoS.Global,
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

	if !c.DisabledHealthCheck {
		go c.healthCheck()
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

	log.Println("Waiting for messages...")

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
func (c *Connection) ListenNotifyClose(done chan bool) {
	connClose := make(chan *amqp.Error)
	c.Conn.NotifyClose(connClose)

	go func() {
		for {
			select {
			case err := <-connClose:
				log.Println("rmq connection lost: ", err)
				log.Printf("reconnecting to rmq in %v...\n", c.ReconnectTime.String())

				c.Retrying = true

				time.Sleep(c.ReconnectTime)

				if err := c.validateHost(); err != nil {
					killService("failed to validate rmq host: ", err)
				}

				if err := c.recreateConn(); err != nil {
					killService("failed to recreate rmq connection: ", err)
				}

				log.Println("sending signal 1 to rmq connection...")

				c.ResetSignal <- 1

				log.Println("signal 1 sent to rmq connection")

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

// recreateConn for rmq
func (c *Connection) recreateConn() error {
	log.Println("trying to recreate rmq connection for host: ", c.Config.Host)

	// important step!
	// prevent healthCheck() to be run once again in c.Setup() because first call to c.Setup() on service Init() already started it (by default)
	// so we do not need/want it to be run again, it would start useless goroutine
	c.DisabledHealthCheck = true

	return c.Setup()
}

// healthCheck for rmq connection
func (c *Connection) healthCheck() {
	for {
		select {
		case <-time.After(healthCheckTime):
			if !c.Retrying {
				// capture current rmq host
				oldHost := c.Config.Host

				if err := c.validateHost(); err != nil {
					killService("failed to validate rmq host: ", err)
				}

				// this means new host was assigned meanwhile (in c.validateHost())
				if oldHost != c.Config.Host {
					if err := c.recreateConn(); err != nil {
						killService("failed to recreate rmq connection: ", err)
					}

					log.Println("rmq connected to new host: ", c.Config.Host)
				}
			}
		}
	}
}

// validateHost will check if rmq host is still valid
// if its invalid -> will resolve dns and assign first valid ip address to rmq host for any further reconnections, c.ConfigHost = <new host ip>
// if its valid still -> nothing happens
func (c *Connection) validateHost() error {
	if check := checkIPConnection(c.Config.Host, c.Config.Port); check {
		return nil
	}

	ips, err := resolveDNS(c.Config.Host)
	if err != nil {
		log.Println("failed to resolve host: ", err)

		return err
	}

	for _, ip := range ips {
		if check := checkIPConnection(ip.String(), c.Config.Port); check {
			c.Config.Host = ip.String()

			break
		}
	}

	return nil
}

// killService with message passed to console output
func killService(msg ...interface{}) {
	log.Println(msg...)
	os.Exit(101)
}

// checkIPConnection will check if IP is available
func checkIPConnection(host string, port string) bool {
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		return false
	}
	defer conn.Close()

	return true
}

// resolveDNS will return assigned ip addresses to given host/record
func resolveDNS(record string) ([]net.IP, error) {
	ips, err := net.LookupIP(record)
	if err != nil {
		return nil, err
	}

	return ips, nil
}
