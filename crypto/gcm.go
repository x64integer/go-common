package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"strings"

	"github.com/semirm-dev/go-dev/util"
)

// GCM crypter
type GCM struct {
	Secret string
}

// Encrypt payload using AES GCM encryption mode
func (gcmEnc *GCM) Encrypt(input []byte) (string, string, string, error) {
	if strings.TrimSpace(gcmEnc.Secret) == "" {
		return "", "", "", errors.New("secret key not provided")
	}

	key := []byte(gcmEnc.Secret)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", "", err
	}

	gcm, err := cipher.NewGCM(block)

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", "", "", err
	}

	encrypted := gcm.Seal(nonce, nonce, input, nil)

	return string(encrypted), hex.EncodeToString(encrypted), util.Base64URLEncode(string(encrypted)), nil
}

// Decrypt AES GCM encrypted input
func (gcmEnc *GCM) Decrypt(input string) (string, error) {
	if strings.TrimSpace(gcmEnc.Secret) == "" {
		return "", errors.New("secret key not provided")
	}

	key := []byte(gcmEnc.Secret)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	byteIn := []byte(input)
	nonceSize := gcm.NonceSize()

	if len(byteIn) < nonceSize {
		return "", errors.New("encrypted text too short")
	}

	nonce, encrypted := byteIn[:nonceSize], byteIn[nonceSize:]

	decrypted, err := gcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}
