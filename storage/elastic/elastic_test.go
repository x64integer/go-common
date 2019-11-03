package elastic_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/semirm-dev/go-dev/storage/elastic"
	"github.com/semirm-dev/go-dev/util"
)

func TestNewConfig(t *testing.T) {
	expected := &elastic.Config{
		Host: util.Env("ELASTIC_HOST", "127.0.0.1"),
		Port: util.Env("ELASTIC_PORT", "9200"),
	}

	config := elastic.NewConfig()

	assert.Equal(t, expected, config)
}
