package sql_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/semirm-dev/go-common/storage/sql"
	"github.com/semirm-dev/go-common/util"
)

func TestNewConfig(t *testing.T) {
	expected := &sql.Config{
		Driver:   util.Env("SQL_DRIVER", "postgres"),
		Host:     util.Env("SQL_HOST", "localhost"),
		Name:     util.Env("SQL_NAME", ""),
		Port:     util.Env("SQL_PORT", "5432"),
		User:     util.Env("SQL_USER", "postgres"),
		Password: util.Env("SQL_PASSWORD", "postgres"),
		SSLMode:  util.Env("SSLMODE", "disable"),
		MaxConn:  20,
	}

	config := sql.NewConfig()

	assert.Equal(t, expected, config)
}
