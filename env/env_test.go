package env_test

import (
	"testing"

	"github.com/semirm-dev/go-dev/env"
	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	gopath := env.Get("GOPATH", "")

	assert.NotEmpty(t, gopath, "Env should not be empty. Check existence of GOPATH")

	non := env.Get("NOT_EXISTS_ENV", "default")

	assert.Equal(t, "default", non, "Env should equal to <default>")
}
