package redis

import (
	"errors"

	"github.com/go-redis/redis"
)

// Storage struct to work with Redis
type Storage struct {
	*Config
}

// Init will create redis client
func (s *Storage) Init() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     s.Config.Host + ":" + s.Config.Port,
		Password: s.Config.Password, // no password set
		DB:       s.Config.DB,       // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, errors.New("[Redis-Ping-Result] - " + err.Error())
	}

	return client, nil
}
