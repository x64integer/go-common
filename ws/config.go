package ws

import "github.com/x64puzzle/go-common/util"

// Config struct
type Config struct {
	WSURL string
}

// NewConfig will init config struct
func NewConfig() *Config {
	return &Config{
		WSURL: util.Env("WS_URL", ""),
	}
}
