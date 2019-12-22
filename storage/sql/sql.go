package sql

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/semirm-dev/go-dev/env"
)

const (
	// PostgresDriver ...
	PostgresDriver = "postgres"
	// MySQLDriver ...
	MySQLDriver = "mysql"
	// MSSQLDriver ...
	MSSQLDriver = "mssql"
)

// Connection for SQL
type Connection struct {
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
	config := &Config{
		Driver:   env.Get("SQL_DRIVER", PostgresDriver),
		Host:     env.Get("SQL_HOST", "localhost"),
		Port:     env.Get("SQL_PORT", "5432"),
		Name:     env.Get("SQL_NAME", ""),
		User:     env.Get("SQL_USER", "postgres"),
		Password: env.Get("SQL_PASSWORD", "postgres"),
		SSLMode:  env.Get("SSLMODE", "disable"),
	}

	maxConn, err := strconv.Atoi(env.Get("SQL_MAX_DB_CONN", "20"))
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
		logrus.Fatal("missing database name")
	}

	var connString string

	driver := strings.ToLower(sqlConn.Config.Driver)

	switch driver {
	case PostgresDriver:
		connString = fmt.Sprintf(
			"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
			sqlConn.Config.Host,
			sqlConn.Config.Port,
			sqlConn.Config.Name,
			sqlConn.Config.User,
			sqlConn.Config.Password,
			sqlConn.Config.SSLMode,
		)
	case MySQLDriver:
		connString = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			sqlConn.Config.User,
			sqlConn.Config.Password,
			sqlConn.Config.Host,
			sqlConn.Config.Port,
			sqlConn.Config.Name,
		)
	case MSSQLDriver:
		connString = fmt.Sprintf(
			"server=%s;port=%s;database=%s;user id=%s;password=%s;",
			sqlConn.Config.Host,
			sqlConn.Config.Port,
			sqlConn.Config.Name,
			sqlConn.Config.User,
			sqlConn.Config.Password,
		)

	default:
		logrus.Fatal("no such SQL driver: ", driver)
	}

	return connString
}
