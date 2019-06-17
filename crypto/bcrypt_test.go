package crypto_test

import (
	"testing"

	"github.com/semirm-dev/go-common/crypto"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var (
	plainBCrypt  = "pwd-123"
	hashedBCrypt = "$2a$10$mYrJ.t1kHdKXYEXI/E2xnunexoJ1mkFJm4PXDWBqhDBbB5TPI0hsm"
)

func TestNewBCrypt(t *testing.T) {
	expected := &crypto.BCrypt{
		Cost: bcrypt.DefaultCost,
	}

	bCrypt := crypto.NewBCrypt()

	assert.Equal(t, expected, bCrypt)
}

func TestBCryptHash(t *testing.T) {
	bCrypt := crypto.NewBCrypt()

	bCrypt.Plain = plainBCrypt

	err := bCrypt.Hash()

	assert.NoError(t, err)
	assert.NotEmpty(t, bCrypt.Hashed)
}

func TestBCryptValidate(t *testing.T) {
	bCrypt := crypto.NewBCrypt()

	bCrypt.Plain = plainBCrypt
	bCrypt.Hashed = hashedBCrypt

	valid := bCrypt.Validate()

	assert.True(t, valid)
}

func TestBCryptHashPlainMissing(t *testing.T) {
	bCrypt := crypto.NewBCrypt()

	bCrypt.Plain = ""

	err := bCrypt.Hash()

	assert.Equal(t, crypto.ErrMissingPlain, err)
}
