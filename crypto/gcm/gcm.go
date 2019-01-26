package gcm

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"strings"

	"github.com/x64integer/go-common/util"
)

// Encrypt will encrypt given input using AES GCM encryption mode
// will return original encrypted value, hex and base64 encoded versions
func Encrypt(input string) (string, string, string, error) {
	secret := util.Env("CRYPTO_SECRET", "")
	if strings.TrimSpace(secret) == "" {
		return "", "", "", errors.New("missing CRYPTO_SECRET env value")
	}

	key := []byte(secret)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", "", err
	}

	gcm, err := cipher.NewGCM(block)

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", "", "", err
	}

	byteIn := []byte(input)
	encrypted := gcm.Seal(nonce, nonce, byteIn, nil)

	return string(encrypted), hex.EncodeToString(encrypted), util.Base64URLEncode(string(encrypted)), nil
}

// Decrypt will decrypt given AES GCM encrypted input
func Decrypt(input string) (string, error) {
	secret := util.Env("CRYPTO_SECRET", "")
	if strings.TrimSpace(secret) == "" {
		return "", errors.New("missing CRYPTO_SECRET env value")
	}

	key := []byte(secret)

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
