package storage

import (
	"log"

	"github.com/x64integer/go-common/storage/cache"
	"github.com/x64integer/go-common/storage/elastic"

	"github.com/x64integer/go-common/storage/redis"
	"github.com/x64integer/go-common/storage/sql"
)

const (
	SQLClient     = 1
	RedisClient   = 2
	ElasticClient = 4
	CacheClient   = 8
)

// Container with storage clients/instances
type Container struct {
	Clients int
	SQL     *sql.Connection
	Redis   *redis.Connection
	Elastic *elastic.Connection
	Cache   cache.Service
}

// Connect and initialize all clients in storage container
func (cont *Container) Connect() {
	if cont.SQL != nil {
		if err := cont.SQL.Connect(); err != nil {
			log.Fatalln("sql connection failed: ", err)
		}
	}

	if cont.Redis != nil {
		if err := cont.Redis.Initialize(); err != nil {
			log.Fatalln("redis initialization failed: ", err)
		}
	}

	if cont.Elastic != nil {
		if err := cont.Elastic.Initialize(); err != nil {
			log.Fatalln("elasticsearch initialization failed: ", err)
		}
	}
}

// DefaultContainer will initialize default storage container
func DefaultContainer(flag int) *Container {
	container := &Container{}

	if flag&SQLClient != 0 {
		container.SQL = &sql.Connection{
			Config: sql.NewConfig(),
		}
	}

	if flag&RedisClient != 0 {
		redisConn := &redis.Connection{
			Config: redis.NewConfig(),
		}

		redisCacheAdapter := &redis.CacheAdapter{
			Connection: redisConn,
		}

		container.Redis = redisConn
		container.Cache = redisCacheAdapter
	}

	if flag&ElasticClient != 0 {
		container.Elastic = &elastic.Connection{
			Config: elastic.NewConfig(),
		}
	}

	return container
}
