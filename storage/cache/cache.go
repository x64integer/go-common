package cache

import "time"

// Service to cache Item
type Service interface {
	// Store item(s) into cache
	Store(...*Item) error
	// Get item(s) from cache
	Get(...*Item) ([]byte, error)
	// Delete item(s) from cache
	Delete(...*Item) error
	// Truncate all items from cache
	Truncate() error
	// Custom func to run against item(s)
	Custom(func(...*Item) error, ...*Item) error
}

// Item to store in Cache
type Item struct {
	Key        string
	Value      interface{}
	Expiration time.Duration
}
