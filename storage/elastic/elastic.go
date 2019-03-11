package elastic

import (
	"github.com/olivere/elastic"
	"github.com/x64integer/go-common/util"
)

// Connection for ElasticSearch
type Connection struct {
	*elastic.Client
	*Config
}

// Config for Elasticsearch connection
type Config struct {
	Host  string
	Port  string
	Sniff bool
}

// NewConfig will initialize default config struct for Elasticsearch
func NewConfig() *Config {
	return &Config{
		Host: util.Env("ELASTIC_HOST", "127.0.0.1"),
		Port: util.Env("ELASTIC_PORT", "9200"),
	}
}

// Initialize Elasticsearch client
func (conn *Connection) Initialize() error {
	client, err := elastic.NewClient(
		elastic.SetURL("http://"+conn.Config.Host+":"+conn.Config.Port),
		elastic.SetSniff(conn.Config.Sniff),
	)

	conn.Client = client

	return err
}
