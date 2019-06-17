package crypto_test

import (
	"testing"

	"github.com/semirm-dev/go-common/crypto"
	"github.com/stretchr/testify/assert"
)

var (
	plainSCrypt  = "pwd-123"
	hashedSCrypt = "32768$8$1$724ee08841de8b8aa8473bbba26250b460ab5f20a658c19709a433ea0f048db1$cae028965d1ccd0015a75d4b33486bf100b15765c3327d3f1ce7a8ba83819ce6"
)

func TestNewSCrypt(t *testing.T) {
	expected := &crypto.SCrypt{
		N:       32768,
		R:       8,
		P:       1,
		KeyLen:  32,
		SaltLen: 32,
	}

	sCrypt := crypto.NewSCrypt()

	assert.Equal(t, expected, sCrypt)
}

func TestSCryptHash(t *testing.T) {
	sCrypt := crypto.NewSCrypt()

	sCrypt.Plain = plainSCrypt

	err := sCrypt.Hash()

	assert.NoError(t, err)
	assert.NotEmpty(t, sCrypt.Salt)
	assert.NotEmpty(t, sCrypt.Hashed)
}

func TestSCryptValidate(t *testing.T) {
	sCrypt := crypto.NewSCrypt()

	sCrypt.Plain = plainSCrypt
	sCrypt.Hashed = hashedSCrypt

	valid := sCrypt.Validate()

	assert.True(t, valid)
}

func TestSCryptHashPlainMissing(t *testing.T) {
	sCrypt := crypto.NewSCrypt()

	sCrypt.Plain = ""

	err := sCrypt.Hash()

	assert.Equal(t, crypto.ErrMissingPlain, err)
}
