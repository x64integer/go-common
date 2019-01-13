package pg

import (
	"database/sql"
)

var (
	// Client is connection to database
	// Exposing public connection variable because connection will be opened one time and used all around
	Client *sql.DB
	// err - error occured during connection initilization
	err error
)

// Storage struct to work with SQL
type Storage struct{}

// InitConnection implements storage.service.InitConnection()
func (s *Storage) InitConnection() error {
	conf := NewConfig() // default values will be used from env variables

	Client, err = s.init(conf)
	if err != nil {
		return err
	}

	return nil
}

// init is helper function to initialize PG connection
func (s *Storage) init(conf *Config) (*sql.DB, error) {
	str := "user=" + conf.User + " password=" + conf.Password + " dbname=" + conf.Name + " sslmode=" + conf.SSLMode

	client, err := sql.Open("postgres", str)
	if err != nil {
		return nil, err
	}

	err = client.Ping()
	if err != nil {
		return nil, err
	}

	client.SetMaxOpenConns(conf.MaxConn)

	return client, nil
}
