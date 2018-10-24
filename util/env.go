package util

import (
	"os"
)

// Env will get os env variable with def as default value
func Env(env, def string) string {
	e := os.Getenv(env)

	if e == "" {
		e = def
	}

	return e
}
