package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-redis/redis"
)

// pipeLength defines limit whether to use pipeline or not
const pipeLength = 1

// Redis cache implementation
type Redis struct {
	Client *redis.Client
}

// Store implements cache.Storage.Store()
func (r *Redis) Store(item ...*Item) error {
	if len(item) > pipeLength { // with pipeline
		pipe := r.Client.Pipeline()

		for _, i := range item {
			pipe.Set(i.Key, i.Value, i.Expiration)
		}

		_, err := pipe.Exec()
		if err != nil {
			return err
		}
	} else { // without pipeline
		var errMsgs []string

		for _, i := range item {
			if err := r.Client.Set(i.Key, i.Value, i.Expiration).Err(); err != nil {
				errMsgs = append(errMsgs, err.Error())
			}
		}

		if len(errMsgs) > 0 {
			return errors.New(strings.Join(errMsgs, ","))
		}
	}

	return nil
}

// Get implements cache.Storage.Get()
func (r *Redis) Get(item ...*Item) ([]byte, error) {
	var result []byte

	if len(item) > pipeLength { // with pipeline

		pipe := r.Client.Pipeline()

		for _, i := range item {
			pipe.Get(i.Key)
		}

		res, err := pipe.Exec()
		if err != nil {
			return nil, err
		}

		var items [][]byte
		for _, r := range res {
			items = append(items, []byte(r.(*redis.StringCmd).Val()))
		}

		itemsByte, err := json.Marshal(items)
		if err != nil {
			return nil, err
		}

		result = itemsByte
	} else { // without pipeline
		var errMsgs []string

		for _, i := range item {
			val, err := r.Client.Get(i.Key).Result()

			switch {
			// key does not exist
			case err == redis.Nil:
				errMsgs = append(errMsgs, fmt.Sprintf("key %v does not exist", i.Key))
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

// Delete implements cache.Storage.Delete()
func (r *Redis) Delete(item ...*Item) error {
	var keys []string

	for _, i := range item {
		keys = append(keys, i.Key)
	}

	return r.Client.Del(keys...).Err()
}

// Truncate implements cache.Storage.Truncate()
func (r *Redis) Truncate() error {
	return r.Client.FlushAll().Err()
}

// Custom implements cache.Storage.Custom()
func (r *Redis) Custom(fn func(...*Item) error, item ...*Item) error {
	return fn(item...)
}
