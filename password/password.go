package password

import (
	"crypto/rand"
	"io"

	"golang.org/x/crypto/bcrypt"
)

// Hash will generate bcrypt-ed password
func Hash(pwd string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// Valid will check hashed password against plain pwd
func Valid(hashed, pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pwd))

	return err == nil
}

// GenerateSalt with given length
// 32 or 64 in most cases
func GenerateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)

	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}

	return salt, nil
}
