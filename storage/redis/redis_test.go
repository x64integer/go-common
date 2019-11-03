package redis_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/semirm-dev/go-dev/storage/redis"
	"github.com/semirm-dev/go-dev/util"
)

func TestNewConfig(t *testing.T) {
	expected := &redis.Config{
		Host:     util.Env("REDIS_HOST", ""),
		Port:     util.Env("REDIS_PORT", "6379"),
		Password: util.Env("REDIS_PASSWORD", ""),
		DB:       0,
	}

	config := redis.NewConfig()

	assert.Equal(t, expected, config)
}
