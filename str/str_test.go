package str_test

import (
	"testing"

	"github.com/semirm-dev/godev/str"

	"github.com/stretchr/testify/assert"
)

func TestUUID(t *testing.T) {
	str := str.UUID()

	assert.NotEmpty(t, str, "Returned UUID must not be empty")
}

func TestRandom(t *testing.T) {
	iter := []int{5, 16, 32}

	for _, i := range iter {
		str := str.Random(i)

		assert.Equal(t, i, len(str))
	}
}
