package cassandra

import (
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/x64integer/go-common/util"
)

// Connection for cassandra
type Connection struct {
	*Config
	Cluster *gocql.ClusterConfig
	Session *gocql.Session
}

// Config for cassandra connection
type Config struct {
	Hosts        []string
	Keyspace     string
	Username     string
	Password     string
	Timeout      time.Duration
	ProtoVersion int
}

// Iterator for select query
type Iterator interface {
	Scan(...interface{}) bool
	MapScan(map[string]interface{}) bool
}

// NewConfig will initialize default config struct for cassandra
func NewConfig() *Config {
	var hosts []string

	hostsEnv := util.Env("CASSANDRA_HOSTS", "127.0.0.1")
	for _, host := range strings.Split(hostsEnv, ",") {
		hosts = append(hosts, host)
	}

	return &Config{
		Hosts:        hosts,
		Keyspace:     util.Env("CASSANDRA_KEYSPACE", "default_keyspace"),
		Username:     util.Env("CASSANDRA_USERNAME", ""),
		Password:     util.Env("CASSANDRA_PASSWORD", ""),
		Timeout:      5 * time.Second,
		ProtoVersion: 4,
	}
}

// Initialize cassandra connection
func (conn *Connection) Initialize() error {
	conn.initCluster()

	if err := conn.NewSession(); err != nil {
		return err
	}

	return nil
}

// NewSession will initialize cassandra session
func (conn *Connection) NewSession() error {
	if conn.Cluster == nil {
		conn.initCluster()
	}

	session, err := conn.Cluster.CreateSession()
	if err != nil {
		return err
	}

	conn.Session = session

	return nil
}

// Close cassandra session
func (conn *Connection) Close() {
	conn.Session.Close()
}

// Exec query against cassandra, non-return query (INSERT, UDPATE, DELETE)
func (conn *Connection) Exec(stmt string, params ...interface{}) error {
	err := conn.Session.Query(stmt, params...).Exec()

	return err
}

// Select data from cassandra
func (conn *Connection) Select(stmt string, params ...interface{}) Iterator {
	iterator := conn.Session.Query(stmt, params...).Iter()

	return iterator
}

// initCluster is helper function to initialize gocql.Cluster
func (conn *Connection) initCluster() {
	cluster := gocql.NewCluster(conn.Config.Hosts...)
	cluster.Keyspace = conn.Config.Keyspace

	// if required from cluster setup
	if conn.Config.Username != "" && conn.Config.Password != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: conn.Config.Username,
			Password: conn.Config.Password,
		}
	}

	cluster.Timeout = conn.Config.Timeout
	cluster.ProtoVersion = conn.Config.ProtoVersion

	conn.Cluster = cluster
}
