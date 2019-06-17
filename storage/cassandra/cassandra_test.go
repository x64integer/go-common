package cassandra_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/semirm-dev/go-common/storage/cassandra"
	"github.com/semirm-dev/go-common/util"
)

func TestNewConfig(t *testing.T) {
	expected := &cassandra.Config{
		Hosts:        []string{"127.0.0.1"},
		Keyspace:     util.Env("CASSANDRA_KEYSPACE", "default_keyspace"),
		Username:     util.Env("CASSANDRA_USERNAME", ""),
		Password:     util.Env("CASSANDRA_PASSWORD", ""),
		Timeout:      5 * time.Second,
		ProtoVersion: 4,
	}

	config := cassandra.NewConfig()

	assert.Equal(t, expected, config)
}
