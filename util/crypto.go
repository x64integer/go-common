package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword will generate bcrypt-ed password
func HashPassword(pwd string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// ValidPassword will check hashed password against plain pwd
func ValidPassword(hashed, pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pwd))

	return err == nil
}

// EncryptGCM will encrypt given input using AES GCM encryption mode
// will return original encrypted value, hex and base64 encoded versions
func EncryptGCM(input string) (string, string, string, error) {
	secret := Env("CRYPTO_SECRET", "")
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

	return string(encrypted), hex.EncodeToString(encrypted), Base64URLEncode(string(encrypted)), nil
}

// DecryptGCM will decrypt given AES GCM encrypted input
func DecryptGCM(input string) (string, error) {
	secret := Env("CRYPTO_SECRET", "")
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

// EncryptCBC will encrypt given input using AES encryption CBC mode
// will return original encrypted value, hex and base64 encoded versions
func EncryptCBC(input string) (string, string, string, error) {
	secret := Env("CRYPTO_SECRET", "")
	if strings.TrimSpace(secret) == "" {
		return "", "", "", errors.New("missing CRYPTO_SECRET env value")
	}

	iv := Env("CRYPTO_IV", "")

	key := []byte(secret)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", "", err
	}

	byteIn := []byte(input)
	byteIn = pkcsPad(byteIn, aes.BlockSize)

	encrypted := make([]byte, len(byteIn))

	byteIV := []byte(iv)

	mode := cipher.NewCBCEncrypter(block, byteIV)
	mode.CryptBlocks(encrypted, byteIn)

	return string(encrypted), hex.EncodeToString(encrypted), Base64URLEncode(string(encrypted)), nil
}

// DecryptCBC will decrypt given AES CBC encrypted input
func DecryptCBC(input string) (string, error) {
	secret := Env("CRYPTO_SECRET", "")
	if strings.TrimSpace(secret) == "" {
		return "", errors.New("missing CRYPTO_SECRET env value")
	}

	iv := Env("CRYPTO_IV", "")

	key := []byte(secret)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	byteIn := []byte(input)
	if len(byteIn) < aes.BlockSize {
		return "", errors.New("encrypted text too short")
	}

	decrypted := make([]byte, len(byteIn))

	byteIV := []byte(iv)

	mode := cipher.NewCBCDecrypter(block, byteIV)
	mode.CryptBlocks(decrypted, byteIn)

	decrypted, err = pkcsUnpad(decrypted, aes.BlockSize)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

// pkcsPad for non-full length blocks
// pkcs5 or pkcs7 will be used based on blocksize
func pkcsPad(ciphertext []byte, blockSize int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(ciphertext, padtext...)
}

// pkcsUnpad will remove PKCS5 padding
// pkcs5 or pkcs7 will be used based on blocksize
func pkcsUnpad(input []byte, blockSize int) ([]byte, error) {
	inputLen := len(input)
	if inputLen == 0 {
		return nil, errors.New("cryptgo/padding: invalid padding size")
	}

	pad := input[inputLen-1]
	padLen := int(pad)
	if padLen > inputLen || padLen > blockSize {
		return nil, errors.New("cryptgo/padding: invalid padding size")
	}

	for _, v := range input[inputLen-padLen : inputLen-1] {
		if v != pad {
			return nil, errors.New("cryptgo/padding: invalid padding")
		}
	}

	return input[:inputLen-padLen], nil
}
