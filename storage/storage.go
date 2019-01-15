package storage

import (
	"database/sql"

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
	// RedisFlag config bit mask
	RedisFlag = 1
	// ElasticFlag config bit mask
	ElasticFlag = 2
	// PgFlag config bit mask
	PgFlag = 4
	// CacheFlag flag
	CacheFlag = 8
)

var (
	// Flag config exposed
	Flag int
	// Redis client exposed
	Redis *redis.Client
	// PG client exposed
	PG *sql.DB
	// Elastic client exposed
	Elastic *elastic.Client
	// Cache instance exposed
	Cache cache.Storage
)

// Init will initialize storage engine based on given config bit mask
//
// Usage ex:
// - Init(storage.RedisFlag | storage.ElasticFlag | storage.PgFlag) -> will initialize Redis, ElasticSearch and Postgres clients
// - Init(storage.ElasticFlag) -> will initialize ElasticSearch client only
func Init(flag int) error {
	Flag = flag

	if Flag&RedisFlag != 0 {
		if err := initRedis(); err != nil {
			return err
		}
	}

	if Flag&PgFlag != 0 {
		if err := initPg(); err != nil {
			return err
		}
	}

	if Flag&ElasticFlag != 0 {
		if err := initElastic(); err != nil {
			return err
		}
	}

	if Flag&CacheFlag != 0 {
		if err := initCache(); err != nil {
			return err
		}
	}

	return nil
}

func initRedis() error {
	redisStorage := &_redis.Storage{
		Config: _redis.NewConfig(),
	}

	redisClient, err := redisStorage.Init()
	if err != nil {
		return err
	}

	Redis = redisClient

	return nil
}

func initPg() error {
	pgStorage := &_pg.Storage{
		Config: _pg.NewConfig(),
	}

	pgClient, err := pgStorage.Init()
	if err != nil {
		return err
	}

	PG = pgClient

	return nil
}

func initElastic() error {
	elasticStorage := &_elastic.Storage{
		Config: _elastic.NewConfig(),
	}

	elasticClient, err := elasticStorage.Init()
	if err != nil {
		return err
	}

	Elastic = elasticClient

	return nil
}

func initCache() error {
	c := util.Env("CACHE_CLIENT", "redis")

	switch c {
	default:
		if Redis == nil {
			if err := initRedis(); err != nil {
				return err
			}
		}

		Cache = &cache.Redis{
			Client: Redis,
		}
	}

	return nil
}
