package mqtt_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/semirm-dev/go-dev/env"
	"github.com/semirm-dev/go-dev/mqtt"
	"github.com/semirm-dev/go-dev/str"
)

func TestNewConfig(t *testing.T) {
	expected := &mqtt.Config{
		Host:          env.Get("MQTT_HOST", "localhost"),
		Port:          env.Get("MQTT_PORT", "1883"),
		Username:      env.Get("MQTT_USERNAME", "guest"),
		Password:      env.Get("MQTT_PASSWORD", "guest"),
		ClientID:      env.Get("MQTT_CLIENT_ID", str.UUID()),
		PubQoS:        0,
		SubQoS:        0,
		CleanSession:  true,
		AutoReconnect: true,
		Retained:      false,
		KeepAlive:     15 * time.Second,
		MsgChanDept:   100,
	}

	config := mqtt.NewConfig()
	expected.ClientID = config.ClientID

	assert.Equal(t, expected, config)
}
