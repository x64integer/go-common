package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/semirm-dev/go-common/storage"
)

func TestDefaultContainer(t *testing.T) {
	st := storage.DefaultContainer(storage.SQLClient | storage.RedisClient | storage.ElasticClient | storage.CassandraClient | storage.CacheClient)

	assert.NotEmpty(t, st.Cassandra)
	assert.NotEmpty(t, st.SQL)
	assert.NotEmpty(t, st.Redis)
	assert.NotEmpty(t, st.Elastic)
	assert.NotEmpty(t, st.Cache)
}
