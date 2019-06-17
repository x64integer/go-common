package crypto

import (
	"crypto/rand"
	"io"
)

// GenerateSalt with given length, 32 or 64 in most cases
func GenerateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)

	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}

	return salt, nil
}
