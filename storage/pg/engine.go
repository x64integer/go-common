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

// Engine struct to work with SQL
type Engine struct{}

// InitConnection implements storage.service.InitConnection()
func (e *Engine) InitConnection() error {
	config := NewConfig() // default values will be used from env variables

	Client, err = e.init(config)
	if err != nil {
		return err
	}

	return nil
}

// init is helper function to initialize PG connection
func (e *Engine) init(config *Config) (*sql.DB, error) {
	str := "user=" + config.User + " password=" + config.Password + " dbname=" + config.Name + " sslmode=disable"

	client, err := sql.Open("postgres", str)
	if err != nil {
		return nil, err
	}

	err = client.Ping()
	if err != nil {
		return nil, err
	}

	client.SetMaxOpenConns(config.MaxConn)

	return client, nil
}
