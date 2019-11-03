package redis

import (
	"github.com/semirm-dev/go-dev/storage/cache"
)

// CacheAdapter will make Connection compatible with cache.Service interface
type CacheAdapter struct {
	*Connection
}

// Store implements and adapts cache.Service.Store
func (adapter *CacheAdapter) Store(items ...*cache.Item) error {
	return adapter.Connection.Store(adapter.toRedisItems(items...)...)
}

// Get implements and adapts cache.Service.Get
func (adapter *CacheAdapter) Get(items ...*cache.Item) ([]byte, error) {
	return adapter.Connection.Get(adapter.toRedisItems(items...)...)
}

// Delete implements and adapts cache.Service.Delete
func (adapter *CacheAdapter) Delete(items ...*cache.Item) error {
	return adapter.Connection.Delete(adapter.toRedisItems(items...)...)
}

// Truncate implements and adapts cache.Service.Truncate
func (adapter *CacheAdapter) Truncate() error {
	return adapter.Connection.Truncate()
}

// Custom implements and adapts cache.Service.Custom
//
// TODO: test calls to custom func fn()
func (adapter *CacheAdapter) Custom(fn func(...*cache.Item) error, cacheItems ...*cache.Item) error {
	redisItems := adapter.toRedisItems(cacheItems...)

	customFunc := func(...*Item) error {
		return fn(cacheItems...)
	}

	return adapter.Connection.Custom(customFunc, redisItems...)
}

// toRedisItems will convert cache.Items to Items
func (adapter *CacheAdapter) toRedisItems(items ...*cache.Item) []*Item {
	var redisItems []*Item

	for _, item := range items {
		redisItems = append(redisItems, &Item{
			Key:        item.Key,
			Value:      item.Value,
			Expiration: item.Expiration,
		})
	}

	return redisItems
}
