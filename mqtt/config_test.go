package mqtt_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/semirm-dev/go-dev/mqtt"
	"github.com/semirm-dev/go-dev/util"
)

func TestNewConfig(t *testing.T) {
	expected := &mqtt.Config{
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

	config := mqtt.NewConfig()
	expected.ClientID = config.ClientID

	assert.Equal(t, expected, config)
}
