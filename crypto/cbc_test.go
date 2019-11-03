package crypto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/semirm-dev/go-dev/crypto"
)

var (
	cbcEncrypted    = "\xe4\xe2V\xef1\x17\x1c?R[\\NS\x8dN0"
	cbcEncryptedHex = "e4e256ef31171c3f525b5c4e538d4e30"
	cbcEncryptedB64 = "5OJW7zEXHD9SW1xOU41OMA=="
)

func TestCBCEncrypt(t *testing.T) {
	cipher := &crypto.Cipher{
		Crypter: &crypto.CBC{
			Secret: secret,
			IV:     iv,
		},
	}

	err := cipher.Encrypt([]byte(message))

	assert.NoError(t, err)
	assert.Equal(t, cbcEncrypted, cipher.Encrypted)
	assert.Equal(t, cbcEncryptedHex, cipher.Hex)
	assert.Equal(t, cbcEncryptedB64, cipher.Base64)
}

func TestCBCDecrypt(t *testing.T) {
	cipher := &crypto.Cipher{
		Crypter: &crypto.CBC{
			Secret: secret,
			IV:     iv,
		},
	}

	err := cipher.Decrypt(cbcEncrypted)

	assert.NoError(t, err)
	assert.Equal(t, message, cipher.Decrypted)
}
