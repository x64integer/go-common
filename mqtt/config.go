package mqtt

import (
	"time"

	"github.com/x64puzzle/go-common/util"
)

// Config for MQTT connection
type Config struct {
	Host          string
	Port          string
	Username      string
	Password      string
	ClientID      string
	KeepAlive     time.Duration
	CleanSession  bool
	AutoReconnect bool
	MsgChanDept   uint
	PubQoS        int
	SubQoS        int
}

// NewConfig will initialize MQTT config struct
func NewConfig() *Config {
	return &Config{
		Host:          util.Env("MQTT_HOST", "localhost"),
		Port:          util.Env("MQTT_PORT", "1883"),
		Username:      util.Env("MQTT_USERNAME", "guest"),
		Password:      util.Env("MQTT_PASSWORD", "guest"),
		ClientID:      util.Env("MQTT_CLIENT_ID", util.UUID()),
		KeepAlive:     15 * time.Second,
		CleanSession:  true,
		AutoReconnect: true,
		MsgChanDept:   100,
		PubQoS:        0,
		SubQoS:        0,
	}
}
