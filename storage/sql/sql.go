package sql

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	"github.com/x64integer/go-common/util"
)

const (
	// PostgresDriver ...
	PostgresDriver = "postgres"
	// MySQLDriver ...
	MySQLDriver = "mysql"
)

// Connection for SQL
type Connection struct {
	Driver string
	*Config
	*sql.DB
}

// Config for SQL connection
type Config struct {
	Driver     string // postgres, mysql
	Host       string // localhost
	Port       string // 5432, 3306
	Name       string // my_db_name
	User       string // my_db_user
	Password   string // my_db_password
	SSLMode    string // disable
	MaxConn    int    // 20
	ConnString string
}

// NewConfig will initialize default config for SQL connection
func NewConfig() *Config {
	config := new(Config)

	config.Driver = util.Env("SQL_DRIVER", PostgresDriver)
	config.Host = util.Env("SQL_HOST", "localhost")
	config.Name = util.Env("SQL_NAME", "")
	config.Port = util.Env("SQL_PORT", "5432")
	config.User = util.Env("SQL_USER", "postgres")
	config.Password = util.Env("SQL_PASSWORD", "postgres")
	config.SSLMode = util.Env("SSLMODE", "disable")

	maxConn, err := strconv.Atoi(util.Env("SQL_MAX_DB_CONN", "20"))
	if err != nil {
		maxConn = 20
	}

	config.MaxConn = maxConn

	return config
}

// Connect to SQL server, open connection
func (sqlConn *Connection) Connect() error {
	if sqlConn.Config.ConnString == "" {
		sqlConn.Config.ConnString = sqlConn.dsn()
	}

	db, err := sql.Open(sqlConn.Config.Driver, sqlConn.Config.ConnString)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(sqlConn.Config.MaxConn)

	if err := db.Ping(); err != nil {
		return err
	}

	sqlConn.DB = db

	return nil
}

// WithDSN will apply connection string for SQL
// If not provided default implementation from *Connection.dsn() will be used
func (sqlConn *Connection) WithDSN(connString string) *Connection {
	sqlConn.Config.ConnString = connString

	return sqlConn
}

// dsn is helper function to construct SQL connection string
func (sqlConn *Connection) dsn() string {
	if sqlConn.Config.Name == "" {
		log.Fatalln("missing database name")
	}

	var connString string

	driver := strings.ToLower(sqlConn.Config.Driver)

	switch driver {
	case PostgresDriver:
		connString = "user=" + sqlConn.Config.User + " password=" + sqlConn.Config.Password + " dbname=" + sqlConn.Config.Name + " sslmode=" + sqlConn.Config.SSLMode
	case MySQLDriver:
		connString = sqlConn.Config.User + ":" + sqlConn.Config.Password + "@tcp(" + sqlConn.Config.Host + ":" + sqlConn.Port + ")" + "/" + sqlConn.Config.Name
	default:
		log.Fatalln("no such SQL driver: " + driver)
	}

	return connString
}
