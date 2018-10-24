// +build unit

package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x64puzzle/go-common/util"
)

func TestJWTGenerate(t *testing.T) {
	jwt := &util.Token{}

	token, err := jwt.Generate(map[string]string{
		"username": "semir",
		"email":    "semir@mail.com",
		"id":       "semir-123",
	})

	assert.NoError(t, err, "Err occured: ", err)
	assert.NotEmpty(t, token, "Token should not be empty")
}

func TestJWTValid(t *testing.T) {
	jwt := &util.Token{}

	token, err := jwt.Generate(map[string]string{
		"username": "semir",
		"email":    "semir@mail.com",
		"id":       "semir-123",
	})

	assert.NoError(t, err, "Err occured: ", err)
	assert.NotEmpty(t, token, "Token should not be empty")

	assert.True(t, jwt.Valid(token), "Token should be valid")
}
