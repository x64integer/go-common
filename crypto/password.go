package crypto

import (
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/scrypt"
)

// Credits to: https://github.com/elithrar/simple-scrypt/blob/master/scrypt.go

// Password hashing
type Password struct {
	Plain   string
	Hashed  string
	Base64  string
	DK      []byte
	Salt    []byte
	N       int // 32768, should be the highest power of 2 derived within 100 milliseconds
	R       int // 8
	P       int // 1
	KeyLen  int // 32
	SaltLen int // 32
}

// NewPassword will initialize default password hash params
func NewPassword() *Password {
	return &Password{
		N:       32768,
		R:       8,
		P:       1,
		KeyLen:  32,
		SaltLen: 32,
	}
}

// Hash password.Plain
func (pwd *Password) Hash() error {
	salt, err := GenerateSalt(pwd.SaltLen)
	if err != nil {
		return err
	}

	pwd.Salt = salt

	dk, err := scrypt.Key([]byte(pwd.Plain), pwd.Salt, pwd.N, pwd.R, pwd.P, pwd.KeyLen)
	if err != nil {
		return err
	}
	pwd.DK = dk

	pwd.Hashed = fmt.Sprintf("%d$%d$%d$%x$%x", pwd.N, pwd.R, pwd.P, pwd.Salt, pwd.DK)
	pwd.Base64 = base64.StdEncoding.EncodeToString(pwd.DK)

	return nil
}

// Validate password.Plain against password.Hashed
func (pwd *Password) Validate() bool {
	// Decode existing hash, retrieve params and salt.
	existing, err := decodeHash(pwd.Hashed)
	if err != nil {
		return false
	}

	// scrypt the pwd.Plain with the same parameters and salt
	b, err := scrypt.Key([]byte(pwd.Plain), existing.Salt, existing.N, existing.R, existing.P, existing.KeyLen)
	if err != nil {
		return false
	}

	// Constant time comparison
	if subtle.ConstantTimeCompare(existing.DK, b) == 1 {
		return true
	}

	return false
}

// decodeHash
func decodeHash(hash string) (*Password, error) {
	pwd := &Password{}

	vals := strings.Split(hash, "$")

	// P, N, R, salt, scrypt derived key
	if len(vals) != 5 {
		return nil, errors.New("invalid hash")
	}

	var err error

	pwd.N, err = strconv.Atoi(vals[0])
	if err != nil {
		return pwd, errors.New("invalid hash")
	}

	pwd.R, err = strconv.Atoi(vals[1])
	if err != nil {
		return pwd, errors.New("invalid hash")
	}

	pwd.P, err = strconv.Atoi(vals[2])
	if err != nil {
		return pwd, errors.New("invalid hash")
	}

	salt, err := hex.DecodeString(vals[3])
	if err != nil {
		return pwd, errors.New("invalid hash")
	}
	pwd.Salt = salt
	pwd.SaltLen = len(salt)

	dk, err := hex.DecodeString(vals[4])
	if err != nil {
		return pwd, errors.New("invalid hash")
	}
	pwd.DK = dk
	pwd.KeyLen = len(dk)

	return pwd, nil
}
