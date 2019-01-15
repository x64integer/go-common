package elastic

import "github.com/x64integer/go-common/util"

// Config for elasticsearch
type Config struct {
	Host  string
	Port  string
	Sniff bool
}

// NewConfig will init config struct
func NewConfig() *Config {
	return &Config{
		Host: util.Env("ELASTIC_HOST", "127.0.0.1"),
		Port: util.Env("ELASTIC_PORT", "9200"),
	}
}
