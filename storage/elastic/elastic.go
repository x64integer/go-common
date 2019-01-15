package elastic

import (
	es "github.com/olivere/elastic"
)

// Storage struct to work with ElasticSearch
type Storage struct {
	*Config
}

// Init will initialize elasticsearch client
func (s *Storage) Init() (*es.Client, error) {
	client, err := es.NewClient(
		es.SetURL("http://"+s.Config.Host+":"+s.Config.Port),
		es.SetSniff(s.Config.Sniff),
	)

	if err != nil {
		return nil, err
	}

	return client, nil
}
