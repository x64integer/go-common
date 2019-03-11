package cache

import "time"

// Service ...
type Service interface {
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

// Item to store in Cache, optional struct
type Item struct {
	Key        string
	Value      interface{}
	Expiration time.Duration
}
