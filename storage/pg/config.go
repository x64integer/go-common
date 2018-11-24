package pg

import (
	"strconv"

	"github.com/x64integer/go-common/util"
)

// Config for postgres
type Config struct {
	Host     string // localhost
	Port     string // 3306, 5432
	Name     string // my_db_name
	User     string // my_db_user
	Password string // my_db_password
	MaxConn  int    // 20
}

// NewConfig will init db connection string
func NewConfig() *Config {
	cs := new(Config)

	cs.Host = util.Env("PG_HOST", "")
	cs.Name = util.Env("PG_NAME", "")
	cs.Port = util.Env("PG_PORT", "5432")
	cs.User = util.Env("PG_USER", "")
	cs.Password = util.Env("PG_PASSWORD", "")

	maxConn, err := strconv.Atoi(util.Env("PG_MAX_DB_CONN", "20"))
	if err != nil {
		maxConn = 20
	}

	cs.MaxConn = maxConn

	return cs
}
