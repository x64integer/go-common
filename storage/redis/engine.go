package redis

import (
	"errors"

	"github.com/go-redis/redis"
)

var (
	// Client is connection to Redis
	Client *redis.Client
	// config
	config = NewConfig()
)

// Engine struct to work with Redis
type Engine struct{}

// InitConnection implements storage.service.InitConnection()
func (e *Engine) InitConnection() error {
	Client = redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password, // no password set
		DB:       0,               // use default DB
	})

	_, err := Client.Ping().Result()
	if err != nil {
		return errors.New("[Redis-Ping-Result] - " + err.Error())
	}

	return nil
}
