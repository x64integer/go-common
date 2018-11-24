package cache

import (
	"time"
)

// Item to query cache Storage interface against
type Item struct {
	Key        string
	Value      interface{}
	Expiration time.Duration
}

// Storage for caching
type Storage interface {
	// Store item/s into cache
	Store(...*Item) error
	// Get items from cache
	Get(...*Item) ([]byte, error)
	// Delete item/s from cache
	Delete(...*Item) error
	// Truncate all stored items in cache
	Truncate() error
	// Custom func to run against item/s
	Custom(func(...*Item) error, ...*Item) error
}
