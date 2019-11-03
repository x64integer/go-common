package jwt_test

import (
	"testing"
	"time"

	"github.com/semirm-dev/go-dev/jwt"
	"github.com/stretchr/testify/assert"
)

var (
	secret = []byte("s4mVxi0fCPYlo1dh1sEWSr4bOWc00krOR1VJ3dSM70BNbW0CTWSbjc2F4tKIzV7h")
	fields = map[string]interface{}{
		"username": "semir",
		"email":    "semir@mail.com",
		"id":       "semir-123",
	}
)

func TestGenerate(t *testing.T) {
	token := &jwt.Token{
		Secret: secret,
	}

	err := token.Generate(&jwt.Claims{
		Expiration: time.Hour * 24,
		Fields:     fields,
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, token.Content)
}

func TestValidateAndExtract(t *testing.T) {
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmF0aW9uIjo4NjQwMDAwMDAwMDAwMCwiRmllbGRzIjp7ImVtYWlsIjoic2VtaXJAbWFpbC5jb20iLCJpZCI6InNlbWlyLTEyMyIsInVzZXJuYW1lIjoic2VtaXIifSwiZXhwIjoxNTYwNTk1OTc3fQ.RWI9Rf3XAVznGdQY3QZ4g7PcNBBRNk02ViWsyqswFVw"

	token := &jwt.Token{
		Secret: secret,
	}

	claims, _ := token.ValidateAndExtract(tokenStr)

	// assert.True(t, valid)
	// assert.NotNil(t, claims)

	assert.NotEmpty(t, claims.Fields)
	assert.Equal(t, claims.Fields, fields)
}

func TestExpiration(t *testing.T) {
	token := &jwt.Token{
		Secret: secret,
	}

	err := token.Generate(&jwt.Claims{
		Expiration: time.Second * 1,
		Fields:     fields,
	})

	assert.NoError(t, err)

	_, valid := token.ValidateAndExtract(token.Content)

	assert.True(t, valid, "token should be valid still")

	time.Sleep(time.Millisecond * 1100)

	_, valid = token.ValidateAndExtract(token.Content)

	assert.True(t, !valid, "token should not be valid")
}

func TestGenerateMissingSecret(t *testing.T) {
	type suite struct {
		Secret  []byte
		Message string
	}

	cases := []*suite{
		&suite{
			Secret:  nil,
			Message: "generate token should fail when secret key is missing",
		},
		&suite{
			Secret:  []byte(""),
			Message: "generate token should fail when secret key length is 0",
		},
		&suite{
			Secret:  []byte{},
			Message: "generate token should fail when secret key length is 0",
		},
	}

	for _, c := range cases {
		token := &jwt.Token{
			Secret: c.Secret,
		}

		err := token.Generate(&jwt.Claims{
			Expiration: time.Hour * 24,
			Fields:     fields,
		})

		assert.Error(t, err, c.Message)
		assert.Equal(t, jwt.ErrMissingSecret, err)
		assert.Empty(t, token.Content)
	}
}
