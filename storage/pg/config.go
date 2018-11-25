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
	SSLMode  string // disable
	MaxConn  int    // 20
}

// NewConfig will init db connection string
func NewConfig() *Config {
	conf := new(Config)

	conf.Host = util.Env("PG_HOST", "")
	conf.Name = util.Env("PG_NAME", "")
	conf.Port = util.Env("PG_PORT", "5432")
	conf.User = util.Env("PG_USER", "")
	conf.Password = util.Env("PG_PASSWORD", "")
	conf.SSLMode = util.Env("SSLMODE", "disable")

	maxConn, err := strconv.Atoi(util.Env("PG_MAX_DB_CONN", "20"))
	if err != nil {
		maxConn = 20
	}

	conf.MaxConn = maxConn

	return conf
}
