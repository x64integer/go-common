package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/semirm-dev/go-common/util"
)

// pipeLength defines limit whether to use pipeline or not
const pipeLength = 1

// Connection for Redis
type Connection struct {
	*redis.Client
	*Config
}

// Config for Redis connection
type Config struct {
	Host       string
	Port       string
	Password   string
	DB         int
	PipeLength int
}

// Item to store in Redis
type Item struct {
	Key        string
	Value      interface{}
	Expiration time.Duration
}

// NewConfig will initialize default config struct for Redis
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

// Initialize Redis client
func (conn *Connection) Initialize() error {
	client := redis.NewClient(&redis.Options{
		Addr:     conn.Config.Host + ":" + conn.Config.Port,
		Password: conn.Config.Password, // no password set
		DB:       conn.Config.DB,       // use default DB
	})

	if conn.Config.PipeLength == 0 {
		conn.Config.PipeLength = pipeLength
	}

	_, err := client.Ping().Result()
	if err != nil {
		return err
	}

	conn.Client = client

	return nil
}

// Store item(s) in Redis
func (conn *Connection) Store(items ...*Item) error {
	if len(items) > conn.PipeLength { // with pipeline
		pipe := conn.Client.Pipeline()

		for _, item := range items {
			pipe.Set(item.Key, item.Value, item.Expiration)
		}

		_, err := pipe.Exec()
		if err != nil {
			return err
		}
	} else { // without pipeline
		var errMsgs []string

		for _, item := range items {
			if err := conn.Client.Set(item.Key, item.Value, item.Expiration).Err(); err != nil {
				errMsgs = append(errMsgs, err.Error())
			}
		}

		if len(errMsgs) > 0 {
			return errors.New(strings.Join(errMsgs, ","))
		}
	}

	return nil
}

// Get item(s) from Redis
func (conn *Connection) Get(items ...*Item) ([]byte, error) {
	var result []byte

	if len(items) > conn.PipeLength { // with pipeline
		pipe := conn.Client.Pipeline()

		for _, item := range items {
			pipe.Get(item.Key)
		}

		res, err := pipe.Exec()
		if err != nil {
			return nil, err
		}

		var itemsToReturn [][]byte
		for _, item := range res {
			itemsToReturn = append(itemsToReturn, []byte(item.(*redis.StringCmd).Val()))
		}

		itemsByte, err := json.Marshal(itemsToReturn)
		if err != nil {
			return nil, err
		}

		result = itemsByte
	} else { // without pipeline
		var errMsgs []string

		for _, item := range items {
			val, err := conn.Client.Get(item.Key).Result()

			switch {
			// key does not exist
			case err == redis.Nil:
				errMsgs = append(errMsgs, fmt.Sprintf("key %v does not exist", item.Key))
			// some other error
			case err != nil:
				errMsgs = append(errMsgs, err.Error())
			// no errors
			case err == nil:
				result = []byte(val)
			}
		}

		if len(errMsgs) > 0 {
			return result, errors.New(strings.Join(errMsgs, ","))
		}
	}

	return result, nil
}

// Delete item(s) from Redis
func (conn *Connection) Delete(items ...*Item) error {
	var keys []string

	for _, item := range items {
		keys = append(keys, item.Key)
	}

	return conn.Client.Del(keys...).Err()
}

// Truncate all items from Redis
func (conn *Connection) Truncate() error {
	return conn.Client.FlushAll().Err()
}

// Custom function to run against item(s)
func (conn *Connection) Custom(fn func(...*Item) error, items ...*Item) error {
	return fn(items...)
}
