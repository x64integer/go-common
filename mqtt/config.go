package mqtt

import (
	"time"

	"github.com/x64integer/go-common/util"
)

// Config for MQTT connection
type Config struct {
	Host          string
	Port          string
	Username      string
	Password      string
	ClientID      string
	PubQoS        int
	SubQoS        int
	CleanSession  bool
	AutoReconnect bool
	Retained      bool
	KeepAlive     time.Duration
	MsgChanDept   uint
}

// NewConfig will initialize MQTT config struct
func NewConfig() *Config {
	return &Config{
		Host:          util.Env("MQTT_HOST", "localhost"),
		Port:          util.Env("MQTT_PORT", "1883"),
		Username:      util.Env("MQTT_USERNAME", "guest"),
		Password:      util.Env("MQTT_PASSWORD", "guest"),
		ClientID:      util.Env("MQTT_CLIENT_ID", util.UUID()),
		PubQoS:        0,
		SubQoS:        0,
		CleanSession:  true,
		AutoReconnect: true,
		Retained:      false,
		KeepAlive:     15 * time.Second,
		MsgChanDept:   100,
	}
}
