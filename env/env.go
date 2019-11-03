package env

import (
	"os"
)

// Get will get os env variable with def as default value
func Get(env, def string) string {
	e := os.Getenv(env)

	if e == "" {
		e = def
	}

	return e
}
