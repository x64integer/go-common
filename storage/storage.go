package storage

import (
	"database/sql"
	"errors"

	"github.com/go-redis/redis"
	"github.com/olivere/elastic"
	"github.com/x64integer/go-common/storage/cache"
	_elastic "github.com/x64integer/go-common/storage/elastic"
	_pg "github.com/x64integer/go-common/storage/pg"
	_redis "github.com/x64integer/go-common/storage/redis"
	"github.com/x64integer/go-common/util"
)

// Caller must close storage connections (as per need)!
// defer storage.PG.Close(), defer storage.Redis.Close(), etc...

const (
	// RedisBitMask config bit mask
	RedisBitMask = 1
	// ElasticBitMask config bit mask
	ElasticBitMask = 2
	// PGBitMask config bit mask
	PGBitMask = 4
	// CacheBitMask flag
	CacheBitMask = 8
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
	// Cache instance exposed
	Cache cache.Storage
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

	// Initialize cache client
	if EngineBitMask&CacheBitMask != 0 {
		c := util.Env("CACHE_CLIENT", "redis")

		switch c {
		default:
			if Redis == nil {
				return errors.New("redis client not initialized - worker should use RedisBitMask in its StorageBitMask()")
			}

			Cache = &cache.Redis{
				Client: Redis,
			}
		}
	}

	return nil
}
