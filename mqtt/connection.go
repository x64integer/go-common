package mqtt

import (
	"errors"

	"github.com/eclipse/paho.mqtt.golang"
)

// Connection struct for MQTT
type Connection struct {
	Config *Config
	Client mqtt.Client
}

// NewConnection will initialize MQTT connection
func NewConnection(config *Config) (*Connection, error) {
	conn := &Connection{
		Config: config,
	}

	opts := mqtt.NewClientOptions()

	broker := conn.Config.Host + ":" + conn.Config.Port

	opts.AddBroker(broker)
	opts.SetClientID(conn.Config.ClientID)
	opts.SetUsername(conn.Config.Username)
	opts.SetPassword(conn.Config.Password)
	opts.SetCleanSession(conn.Config.CleanSession)
	opts.SetAutoReconnect(conn.Config.AutoReconnect)
	opts.SetKeepAlive(conn.Config.KeepAlive)
	opts.SetMessageChannelDepth(conn.Config.MsgChanDept)

	conn.Client = mqtt.NewClient(opts)
	if token := conn.Client.Connect(); token.Wait() && token.Error() != nil {
		return nil, errors.New("MQTT client connection failed: " + token.Error().Error())
	}

	return conn, nil
}

// Publish payload p to topic t
func (conn *Connection) Publish(t string, p []byte) mqtt.Token {
	token := conn.Client.Publish(t, byte(conn.Config.PubQoS), conn.Config.Retained, p)

	return token
}

// Subscribe to topic c
func (conn *Connection) Subscribe(t string, callback func(c mqtt.Client, m mqtt.Message)) mqtt.Token {
	token := conn.Client.Subscribe(t, byte(conn.Config.SubQoS), callback)

	return token
}
