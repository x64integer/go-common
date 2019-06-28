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
		// client will connect to this url
		WSURL: util.Env("WS_URL", ""),
		// server will open connection on this endpoint
		Host:     "localhost",
		Port:     "8080",
		Endpoint: "/",
	}
}
