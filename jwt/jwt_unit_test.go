// +build unit

package jwt_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/x64integer/go-common/jwt"
)

var (
	secret   = []byte("s4mVxi0fCPYlo1dh1sEWSr4bOWc00krOR1VJ3dSM70BNbW0CTWSbjc2F4tKIzV7h")
	tokenStr = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InNlbWlyQG1haWwuY29tIiwiZXhwIjoxNTUyMTMzNDI3LCJpZCI6InNlbWlyLTEyMyIsInVzZXJuYW1lIjoic2VtaXIifQ.bASFJHnwo7G_FpHVldUDXxFeYuGTPJyRZi0N4KBNC2g"
)

func TestGenerate(t *testing.T) {
	token := &jwt.Token{
		Secret: secret,
	}

	err := token.Generate(&jwt.Claims{
		Expiration: time.Hour * 24,
		Fields: map[string]interface{}{
			"username": "semir",
			"email":    "semir@mail.com",
			"id":       "semir-123",
		},
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, token.Content)
}

func TestValidateAndExtract(t *testing.T) {
	token := &jwt.Token{
		Secret: secret,
	}

	claims, valid := token.ValidateAndExtract(tokenStr)
	assert.True(t, valid)
	assert.NotEmpty(t, claims)
}
