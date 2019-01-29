package pg

import (
	"database/sql"
)

// Storage struct to work with SQL
type Storage struct {
	*Config
}

// Init will initialize postgres client
func (s *Storage) Init() (*sql.DB, error) {
	str := "user=" + s.Config.User + " password=" + s.Config.Password + " dbname=" + s.Config.Name + " sslmode=" + s.Config.SSLMode

	client, err := sql.Open("postgres", str)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(); err != nil {
		return nil, err
	}

	client.SetMaxOpenConns(s.Config.MaxConn)

	return client, nil
}
