package redis

import "github.com/x64integer/go-common/util"

// Config for redis
type Config struct {
	Host     string
	Port     string
	Password string
}

// NewConfig will init Config struct
func NewConfig() *Config {
	return &Config{
		Host:     util.Env("REDIS_HOST", ""),
		Port:     util.Env("REDIS_PORT", "6379"),
		Password: util.Env("REDIS_PASSWORD", ""),
	}
}
