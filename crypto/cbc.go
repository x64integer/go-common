package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/semirm-dev/godev/str"
)

// CBC crypter
type CBC struct {
	Secret string
	IV     string
}

// Encrypt payload using AES encryption CBC mode
func (cbcEnc *CBC) Encrypt(payload []byte) (string, string, string, error) {
	if strings.TrimSpace(cbcEnc.Secret) == "" {
		return "", "", "", errors.New("secret key not provided")
	}

	key := []byte(cbcEnc.Secret)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", "", err
	}

	byteIn := pkcsPad(payload, aes.BlockSize)

	encrypted := make([]byte, len(byteIn))

	byteIV := []byte(cbcEnc.IV)

	mode := cipher.NewCBCEncrypter(block, byteIV)
	mode.CryptBlocks(encrypted, byteIn)

	return string(encrypted), hex.EncodeToString(encrypted), str.Base64URLEncode(string(encrypted)), nil
}

// Decrypt AES CBC encrypted input
func (cbcEnc *CBC) Decrypt(encrypted string) (string, error) {
	if strings.TrimSpace(cbcEnc.Secret) == "" {
		return "", errors.New("secret key not provided")
	}

	key := []byte(cbcEnc.Secret)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	byteIn := []byte(encrypted)
	if len(byteIn) < aes.BlockSize {
		return "", errors.New("encrypted text too short")
	}

	decrypted := make([]byte, len(byteIn))

	byteIV := []byte(cbcEnc.IV)

	mode := cipher.NewCBCDecrypter(block, byteIV)
	mode.CryptBlocks(decrypted, byteIn)

	decrypted, err = pkcsUnPad(decrypted, aes.BlockSize)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

// pkcsPad for non-full length blocks
// pkcs5 or pkcs7 will be used based on block size
func pkcsPad(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(ciphertext, padtext...)
}

// pkcsUnPad will remove PKCS5 padding
// pkcs5 or pkcs7 will be used based on block size
func pkcsUnPad(input []byte, blockSize int) ([]byte, error) {
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
