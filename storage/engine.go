package storage

import (
	"database/sql"

	"github.com/go-redis/redis"
	"github.com/olivere/elastic"
	_elastic "github.com/x64puzzle/go-common/storage/elastic"
	_pg "github.com/x64puzzle/go-common/storage/pg"
	_redis "github.com/x64puzzle/go-common/storage/redis"
)

// Caller must close storage connections (as per need)!
// defer storage.PG.Close(), defer storage.Redis.Close(), etc...

// In order to add new storage engine, do the following:
// 1.) Add config bit mask constant: 1, 2, 4, 8, 16, 32, 64, 128...
// 2.) Expose its client var
// 3.) Implement storage.service.InitConnection() for new storage engine (take Redis as an example)
// 4.) Initialize its Client in Init(engine int) func

const (
	// RedisBitMask config bit mask
	RedisBitMask = 1
	// ElasticBitMask config bit mask
	ElasticBitMask = 2
	// PGBitMask config bit mask
	PGBitMask = 4
)

var (
	// EngineBitMask config exposed
	EngineBitMask int
	// Redis client exposed
	Redis *redis.Client
	// Elastic client exposed
	Elastic *elastic.Client
	// PG client exposed
	PG *sql.DB
)

// Init will initialize storage engine based on given config bit mask
//
// Usage ex:
// - Init(storage.RedisBitMask | storage.ElasticBitMask | storage.PGBitMask) -> will initialize Redis, ElasticSearch and SQL clients
// - Init(storage.ElasticBitMask) -> will initialize ElasticSearch client only
func Init(engineBitMask int) error {
	EngineBitMask = engineBitMask

	// Initialize Redis Client
	if EngineBitMask&RedisBitMask != 0 {
		engine := &_redis.Engine{}

		if err := engine.InitConnection(); err != nil {
			return err
		}

		Redis = _redis.Client
	}

	// Initialize ElasticSearch Client
	if EngineBitMask&ElasticBitMask != 0 {
		engine := &_elastic.Engine{}

		if err := engine.InitConnection(); err != nil {
			return err
		}

		Elastic = _elastic.Client
	}

	// Initialize PG SQL Client
	if EngineBitMask&PGBitMask != 0 {
		engine := &_pg.Engine{}

		if err := engine.InitConnection(); err != nil {
			return err
		}

		PG = _pg.Client
	}

	return nil
}
