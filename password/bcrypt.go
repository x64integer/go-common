package password

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrMissingPlainBCrypt error
	ErrMissingPlainBCrypt = errors.New("missing Plain property")
)

// BCrypt hashing algorithm
type BCrypt struct {
	Plain  string
	Hashed string
	Cost   int
}

// NewBCrypt will initialize default bcrypt params
func NewBCrypt() *BCrypt {
	return &BCrypt{
		Cost: bcrypt.DefaultCost,
	}
}

// Hash bCrypt.Plain
func (bCrypt *BCrypt) Hash() error {
	if bCrypt.Plain == "" {
		return ErrMissingPlainBCrypt
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(bCrypt.Plain), bCrypt.Cost)
	if err != nil {
		return err
	}

	bCrypt.Hashed = string(hashed)

	return nil
}

// Validate bCrypt.Plain against bCrypt.Hashed
func (bCrypt *BCrypt) Validate() bool {
	return bcrypt.CompareHashAndPassword([]byte(bCrypt.Hashed), []byte(bCrypt.Plain)) == nil
}
