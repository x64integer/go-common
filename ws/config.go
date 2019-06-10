package ws

import "github.com/semirm-dev/go-common/util"

// Config for websocket client/server
type Config struct {
	WSURL    string
	Host     string
	Port     string
	Endpoint string
}

// NewConfig will initialize default configuration
func NewConfig() *Config {
	return &Config{
		WSURL:    util.Env("WS_URL", ""),
		Host:     "",
		Port:     "8080",
		Endpoint: "/",
	}
}
