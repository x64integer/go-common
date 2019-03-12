package elastic

import (
	"context"
	"encoding/json"

	"github.com/olivere/elastic"
	"github.com/x64integer/go-common/util"
)

const (
	defaultSearchLimit = 50
)

// Connection for ElasticSearch
type Connection struct {
	*elastic.Client
	*Config
}

// Config for Elasticsearch connection
type Config struct {
	Host     string
	Port     string
	Sniff    bool
	UseHTTPS bool
}

// Entity to store in elasticsearch
type Entity struct {
	ID      string
	Content interface{}
}

// SearchEntity ...
type SearchEntity struct {
	Term   string
	Fields []string
	Offset int
	Limit  int
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
	link := "http"

	if conn.UseHTTPS {
		link = "https"
	}

	client, err := elastic.NewClient(
		elastic.SetURL(link+"://"+conn.Config.Host+":"+conn.Config.Port),
		elastic.SetSniff(conn.Config.Sniff),
	)

	conn.Client = client

	return err
}

// Insert data into elasticsearch
func (conn *Connection) Insert(ctx context.Context, index string, t string, entity *Entity) error {
	svc := conn.Index().Index(index).Type(t)

	_, err := svc.Id(entity.ID).BodyJson(entity.Content).Do(ctx)

	return err
}

// BulkInsert data into elasticsearch
func (conn *Connection) BulkInsert(ctx context.Context, index string, t string, entities ...*Entity) error {
	bulk := conn.Bulk().Index(index).Type(t)

	for _, entity := range entities {
		bulk.Add(elastic.NewBulkIndexRequest().Id(entity.ID).Doc(entity.Content))
	}

	_, err := bulk.Do(ctx)

	return err
}

// SearchByTerm data from elasticsearch
func (conn *Connection) SearchByTerm(ctx context.Context, index string, t string, searchEntity *SearchEntity) ([]byte, error) {
	if searchEntity.Limit == 0 {
		searchEntity.Limit = defaultSearchLimit
	}

	query := elastic.NewMultiMatchQuery(searchEntity.Term, searchEntity.Fields...)

	searchResult, err := conn.Search().Index(index).Type(t).From(searchEntity.Offset).Size(searchEntity.Limit).Query(query).Do(ctx)
	if err != nil {
		return nil, err
	}

	var resp []interface{}

	for _, hit := range searchResult.Hits.Hits {
		var item interface{}

		if err := json.Unmarshal(*hit.Source, &item); err != nil {
			continue
		}

		resp = append(resp, item)
	}

	b, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	return b, nil
}
