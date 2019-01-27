package api

// Dao represents data access layer
// Database logic
type Dao struct {
}

// Save user account
func (dao *Dao) Save(stmt string) error {
	// _, err := storage.PG.Exec("INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4)", "", "", "", "", "")

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
