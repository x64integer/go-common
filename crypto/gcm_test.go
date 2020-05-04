package crypto_test

import (
	"testing"

	"github.com/semirm-dev/godev/crypto"
	"github.com/stretchr/testify/assert"
)

func TestGCMEncrypt(t *testing.T) {
	cipher := &crypto.Cipher{
		Crypter: &crypto.GCM{
			Secret: secret,
		},
	}

	err := cipher.Encrypt([]byte(message))

	assert.NoError(t, err)
	assert.NotEmpty(t, cipher.Encrypted)
	assert.NotEmpty(t, cipher.Hex)
	assert.NotEmpty(t, cipher.Base64)
}

func TestGCMDecrypt(t *testing.T) {
	cipher := &crypto.Cipher{
		Crypter: &crypto.GCM{
			Secret: secret,
		},
	}

	err := cipher.Decrypt("\xc9:>z\x9d\xde;l\xe0\xd9n\x16#\xc0M\xcd]0Il\xae\x10\xd8L\x1aTd\xc3vܝQ\xc6>!\xfb\x00\t\xb4\xf0")

	assert.NoError(t, err)
	assert.Equal(t, message, cipher.Decrypted)
}
