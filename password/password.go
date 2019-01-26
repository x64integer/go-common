package password

import "golang.org/x/crypto/bcrypt"

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
