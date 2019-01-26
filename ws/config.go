package ws

import "github.com/x64integer/go-common/util"

// Config struct
type Config struct {
	WSURL    string
	Host     string
	Port     string
	Endpoint string
}

// NewConfig will init config struct
func NewConfig() *Config {
	return &Config{
		WSURL:    util.Env("WS_URL", ""),
		Host:     "",
		Port:     "8080",
		Endpoint: "/",
	}
}
