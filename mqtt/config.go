package mqtt

import (
	"time"

	"github.com/semirm-dev/go-dev/util"

	"github.com/semirm-dev/go-dev/env"
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
		Host:          env.Get("MQTT_HOST", "localhost"),
		Port:          env.Get("MQTT_PORT", "1883"),
		Username:      env.Get("MQTT_USERNAME", "guest"),
		Password:      env.Get("MQTT_PASSWORD", "guest"),
		ClientID:      env.Get("MQTT_CLIENT_ID", util.UUID()),
		PubQoS:        0,
		SubQoS:        0,
		CleanSession:  true,
		AutoReconnect: true,
		Retained:      false,
		KeepAlive:     15 * time.Second,
		MsgChanDept:   100,
	}
}
