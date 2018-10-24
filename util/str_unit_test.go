// +build unit

package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x64puzzle/go-common/util"
)

func TestUUID(t *testing.T) {
	str := util.UUID()

	assert.NotEmpty(t, str, "Returned UUID must not be empty")
}

func TestRandomStr(t *testing.T) {
	iter := []int{5, 16, 32}

	for _, i := range iter {
		str := util.RandomStr(i)

		assert.Equal(t, i, len(str))
	}
}
