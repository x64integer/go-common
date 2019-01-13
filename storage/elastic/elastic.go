package elastic

import (
	es "github.com/olivere/elastic"
)

var (
	// Client is connection to ES
	Client *es.Client
	// err - error occured during connection initilization
	err error
	// config
	config = NewConfig()
)

// Storage struct to work with ElasticSearch
type Storage struct{}

// InitConnection implements storage.service.InitConnection()
func (s *Storage) InitConnection() error {
	Client, err = es.NewClient(
		es.SetURL("http://"+config.Host+":"+config.Port),
		es.SetSniff(false),
	)

	if err != nil {
		return err
	}

	return nil
}
