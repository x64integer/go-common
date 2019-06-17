package password_test

import (
	"testing"

	"github.com/semirm-dev/go-common/password"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var (
	plainBCrypt  = "pwd-123"
	hashedBCrypt = "$2a$10$mYrJ.t1kHdKXYEXI/E2xnunexoJ1mkFJm4PXDWBqhDBbB5TPI0hsm"
)

func TestNewBCrypt(t *testing.T) {
	expected := &password.BCrypt{
		Cost: bcrypt.DefaultCost,
	}

	bCrypt := password.NewBCrypt()

	assert.Equal(t, expected, bCrypt)
}

func TestBCryptHash(t *testing.T) {
	bCrypt := password.NewBCrypt()

	bCrypt.Plain = plainBCrypt

	err := bCrypt.Hash()

	assert.NoError(t, err)
	assert.NotEmpty(t, bCrypt.Hashed)
}

func TestBCryptValidate(t *testing.T) {
	bCrypt := password.NewBCrypt()

	bCrypt.Plain = plainBCrypt
	bCrypt.Hashed = hashedBCrypt

	valid := bCrypt.Validate()

	assert.True(t, valid)
}

func TestBCryptHashPlainMissing(t *testing.T) {
	bCrypt := password.NewBCrypt()

	bCrypt.Plain = ""

	err := bCrypt.Hash()

	assert.Equal(t, password.ErrMissingPlainBCrypt, err)
}
