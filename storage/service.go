package storage

// Service for storage engine
type Service interface {
	InitConnection() error
}
