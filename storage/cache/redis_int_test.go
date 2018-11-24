package cache_test

import (
	"testing"

	"github.com/SwipeStoxGmbH/exchange-services/logger"
	"github.com/SwipeStoxGmbH/exchange-services/storage"
	"github.com/SwipeStoxGmbH/exchange-services/storage/cache"
	"github.com/stretchr/testify/assert"
)

func init() {
	if err := storage.Init(storage.RedisBitMask + storage.CacheBitMask); err != nil {
		panic(err)
	}
}

// storeSingle item in cache without pipeline
func storeSingle() error {
	return storage.Cache.Store(&cache.Item{
		Key:   "test_item_1",
		Value: "value 1",
	})
}

// storeMany items in cache using pipeline
func storeMany() error {
	return storage.Cache.Store(&cache.Item{
		Key:   "test_item_2",
		Value: "value 2",
	}, &cache.Item{
		Key:   "test_item_3",
		Value: "value 3",
	}, &cache.Item{
		Key:   "test_item_4",
		Value: "value 4",
	})
}

func TestCacheInstance(t *testing.T) {
	assert.NotNil(t, storage.Cache)
}

func TestRedisStore(t *testing.T) {
	err := storeSingle()

	assert.NoError(t, err)

	err = storeMany()

	assert.NoError(t, err)
}

func TestRedisGet(t *testing.T) {
	err := storeMany()
	assert.NoError(t, err)

	items, err := storage.Cache.Get(&cache.Item{Key: "test_item_1"})
	assert.NoError(t, err)
	assert.NotEmpty(t, items)
	assert.Equal(t, 1, len(items))

	logger.Log.Info("item: ", items)

	items, err = storage.Cache.Get(&cache.Item{Key: "test_item_1"}, &cache.Item{Key: "test_item_3"}, &cache.Item{Key: "test_item_4"})
	assert.NoError(t, err)
	assert.NotEmpty(t, items)
	assert.Equal(t, 3, len(items))

	logger.Log.Info("items: ", items)
}
