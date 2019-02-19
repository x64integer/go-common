package api

import (
	"github.com/x64integer/go-common/storage"
)

// Dao represents data access layer
// Database logic
type Dao struct {
}

// Save user account
func (dao *Dao) Save(stmt string, data []interface{}) error {
	_, err := storage.PG.Exec(stmt, data...)
	if err != nil {
		return err
	}

	return nil
}

// CreateSession will store user token into cache
func (dao *Dao) CreateSession() error {
	return nil
}

// DestroySession will delete user token from cache
func (dao *Dao) DestroySession() error {
	return nil
}
