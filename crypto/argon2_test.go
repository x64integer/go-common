package crypto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/semirm-dev/go-dev/crypto"
)

var (
	plainArgon  = "pwd-123"
	hashedArgon = "19$65536$3$2$26f00e17ceedfe7029fd504cad743d9c61f3bb1ffa8c4f7fdf02df8a1a2e7b68$21484942af11d6fcc7cf4f4a72a81e84b4ec9a4490f1c86055ea203f4007d9d6"
)

func TestNewArgon2(t *testing.T) {
	expected := &crypto.Argon2{
		Memory:  64 * 1024,
		Time:    3,
		Threads: 2,
		SaltLen: 32,
		KeyLen:  32,
	}

	argon := crypto.NewArgon2()

	assert.Equal(t, expected, argon)
}

func TestArgon2Hash(t *testing.T) {
	argon := crypto.NewArgon2()

	argon.Plain = plainArgon

	err := argon.Hash()

	assert.NoError(t, err)

	assert.NotEmpty(t, argon.Salt)
	assert.NotEmpty(t, argon.DK)
	assert.NotEmpty(t, argon.Hashed)
}

func TestArgon2Validate(t *testing.T) {
	argon := crypto.NewArgon2()

	argon.Plain = plainArgon
	argon.Hashed = hashedArgon

	valid := argon.Validate()

	assert.True(t, valid)
}

func TestArgon2HashPlainMissing(t *testing.T) {
	argon := crypto.NewArgon2()

	argon.Plain = ""

	err := argon.Hash()

	assert.Equal(t, crypto.ErrMissingPlain, err)
}
