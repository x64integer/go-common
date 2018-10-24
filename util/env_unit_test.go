// +build unit

package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x64puzzle/go-common/util"
)

func TestEnv(t *testing.T) {
	gopath := util.Env("GOPATH", "")

	assert.NotEmpty(t, gopath, "Env should not be empty. Check existence of GOPATH")

	non := util.Env("NOT_EXISTS_ENV", "default")

	assert.Equal(t, "default", non, "Env should equal to <default>")
}
