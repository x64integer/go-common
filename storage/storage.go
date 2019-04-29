package storage

import (
	"log"

	"github.com/semirm-dev/go-common/storage/cache"
	"github.com/semirm-dev/go-common/storage/cassandra"
	"github.com/semirm-dev/go-common/storage/elastic"

	"github.com/semirm-dev/go-common/storage/redis"
	"github.com/semirm-dev/go-common/storage/sql"
)

const (
	// SQLClient flag
	SQLClient = 1
	// RedisClient flag
	RedisClient = 2
	// ElasticClient flag
	ElasticClient = 4
	// CassandraClient flag
	CassandraClient = 8
)

// C exposes storage container so it can be accessed globally
var C *Container

// Container with storage clients/instances
type Container struct {
	SQL       *sql.Connection
	Redis     *redis.Connection
	Elastic   *elastic.Connection
	Cassandra *cassandra.Connection
	Cache     cache.Service
	Expose    bool
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

	if cont.Cassandra != nil {
		if err := cont.Cassandra.Initialize(); err != nil {
			log.Fatalln("cassandra initialization failed: ", err)
		}
	}

	if cont.Expose {
		C = cont
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

	if flag&CassandraClient != 0 {
		container.Cassandra = &cassandra.Connection{
			Config: cassandra.NewConfig(),
		}
	}

	return container
}
