package elastic

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/olivere/elastic"
	"github.com/semirm-dev/go-common/util"
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

// SearchEntity is used to query Elasticsearch
type SearchEntity struct {
	Term      string
	Fields    []string
	From      int
	Limit     int
	Sort      string
	SortOrder bool // true -> asc, false -> desc
}

// NewConfig will initialize default configuration for Elasticsearch
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

	searchService := conn.Search().
		Index(index).
		Type(t).
		From(searchEntity.From).
		Size(searchEntity.Limit).
		Query(query)

	if searchEntity.Sort != "" {
		searchService.Sort(searchEntity.Sort, searchEntity.SortOrder)
	}

	if strings.TrimSpace(searchEntity.Term) == "" {
		searchService.Query(elastic.NewMatchAllQuery())
	}

	searchResult, err := searchService.Do(ctx)
	if err != nil {
		return nil, err
	}

	return conn.buildResponse(searchResult.Hits.Hits)
}

func (conn *Connection) buildResponse(hits []*elastic.SearchHit) ([]byte, error) {
	var resp []interface{}

	for _, hit := range hits {
		var item interface{}

		if err := json.Unmarshal(*&hit.Source, &item); err != nil {
			continue
		}

		resp = append(resp, item)
	}

	b, err := json.Marshal(resp)

	return b, err
}
