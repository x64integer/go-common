package util

import "golang.org/x/crypto/bcrypt"

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
