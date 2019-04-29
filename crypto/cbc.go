package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/semirm-dev/go-common/util"
)

// CBC crypter
type CBC struct {
	Secret string
	IV     string
}

// Encrypt payload using AES encryption CBC mode
func (cbcEnc *CBC) Encrypt(input []byte) (string, string, string, error) {
	if strings.TrimSpace(cbcEnc.Secret) == "" {
		return "", "", "", errors.New("secret key not provided")
	}

	key := []byte(cbcEnc.Secret)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", "", err
	}

	byteIn := pkcsPad(input, aes.BlockSize)

	encrypted := make([]byte, len(byteIn))

	byteIV := []byte(cbcEnc.IV)

	mode := cipher.NewCBCEncrypter(block, byteIV)
	mode.CryptBlocks(encrypted, byteIn)

	return string(encrypted), hex.EncodeToString(encrypted), util.Base64URLEncode(string(encrypted)), nil
}

// Decrypt AES CBC encrypted input
func (cbcEnc *CBC) Decrypt(input string) (string, error) {
	if strings.TrimSpace(cbcEnc.Secret) == "" {
		return "", errors.New("secret key not provided")
	}

	key := []byte(cbcEnc.Secret)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	byteIn := []byte(input)
	if len(byteIn) < aes.BlockSize {
		return "", errors.New("encrypted text too short")
	}

	decrypted := make([]byte, len(byteIn))

	byteIV := []byte(cbcEnc.IV)

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
