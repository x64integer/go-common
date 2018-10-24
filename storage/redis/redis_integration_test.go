// +build int

package redis_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x64puzzle/go-common/storage"
	"github.com/x64puzzle/go-common/storage/redis"
)

func TestInitConnection(t *testing.T) {
	err := storage.Init(storage.RedisBitMask)
	assert.NoError(t, err)

	assert.NotNil(t, redis.Client, "Make sure redis is running")
	assert.NotNil(t, storage.Redis, "Make sure storage.Redis is exposed in storage.engine")
}
