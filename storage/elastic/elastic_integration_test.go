// +build int

package elastic_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x64puzzle/go-common/storage"
	"github.com/x64puzzle/go-common/storage/elastic"
)

func TestInitConnection(t *testing.T) {
	err := storage.Init(storage.ElasticBitMask)
	assert.NoError(t, err)

	assert.NotNil(t, elastic.Client, "Make sure elasticsearch is running")
	assert.NotNil(t, storage.Elastic, "Make sure storage.Elastic is exposed in /storage.engine")
}
