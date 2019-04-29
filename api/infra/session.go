package infra

import (
	"time"

	"github.com/semirm-dev/go-common/storage/cache"
)

// Session for user
type Session struct {
	Cache cache.Service
}

// Create new session
func (session *Session) Create(key, val string) error {
	if err := session.Cache.Store(&cache.Item{
		Key:        key,
		Value:      val,
		Expiration: 24 * time.Hour,
	}); err != nil {
		return err
	}

	return nil
}

// Get value from session
func (session *Session) Get(key string) (string, error) {
	val, err := session.Cache.Get(&cache.Item{Key: key})
	if err != nil {
		return "", err
	}

	return string(val), nil
}

// Destroy session
func (session *Session) Destroy(key string) error {
	if err := session.Cache.Delete(&cache.Item{Key: key}); err != nil {
		return err
	}

	return nil
}
