package redis

import (
	"strconv"

	"github.com/x64integer/go-common/util"
)

// Config for redis
type Config struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// NewConfig will init Config struct
func NewConfig() *Config {
	db, err := strconv.Atoi(util.Env("REDIS_DB", "0"))
	if err != nil {
		db = 0
	}

	return &Config{
		Host:     util.Env("REDIS_HOST", ""),
		Port:     util.Env("REDIS_PORT", "6379"),
		Password: util.Env("REDIS_PASSWORD", ""),
		DB:       db,
	}
}
