package crypto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/semirm-dev/godev/crypto"
)

var (
	secret  = "kYp3s6v9y$B&E)H+MbQeThWmZq4t7w!z" // 256bit
	iv      = "D*G-KaPdSgVkYp3s"                 // len(secret) / 2 = 128bit
	message = "test message"
)

var (
	mockEncrypted    = "encrypted-message"
	mockEncryptedHex = "encrypted-hex"
	mockEncryptedB64 = "encrypted-base64"
)

type Mock struct{}

func (mock *Mock) Encrypt(input []byte) (string, string, string, error) {
	return mockEncrypted, mockEncryptedHex, mockEncryptedB64, nil
}

func (mock *Mock) Decrypt(input string) (string, error) {
	return message, nil
}

func TestCryptoEncrypt(t *testing.T) {
	cipher := &crypto.Cipher{
		Crypter: &Mock{},
	}

	err := cipher.Encrypt([]byte(message))

	assert.NoError(t, err)
	assert.Equal(t, mockEncrypted, cipher.Encrypted)
	assert.Equal(t, mockEncryptedHex, cipher.Hex)
	assert.Equal(t, mockEncryptedB64, cipher.Base64)
}

func TestCryptoDecrypt(t *testing.T) {
	cipher := &crypto.Cipher{
		Crypter: &Mock{},
	}

	err := cipher.Decrypt(message)

	assert.NoError(t, err)
	assert.Equal(t, message, cipher.Decrypted)
}

func TestCryptoMissingCrypter(t *testing.T) {
	cipher := &crypto.Cipher{}

	err := cipher.Encrypt([]byte(message))

	assert.Error(t, err)
	assert.Equal(t, crypto.ErrMissingCrypter, err)
}
