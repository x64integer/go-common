package password

import (
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/scrypt"
)

// Credits to: https://github.com/elithrar/simple-scrypt/blob/master/scrypt.go

// SCrypt hashing algorithm
type SCrypt struct {
	Plain   string
	Hashed  string
	DK      []byte
	Salt    []byte
	N       int // 32768, should be the highest power of 2 derived within 100 milliseconds
	R       int // 8
	P       int // 1
	KeyLen  int // 32
	SaltLen int // 32
}

// NewSCrypt will initialize default SCrypt params
func NewSCrypt() *SCrypt {
	return &SCrypt{
		N:       32768,
		R:       8,
		P:       1,
		KeyLen:  32,
		SaltLen: 32,
	}
}

// Hash sCrypt.Plain
func (sCrypt *SCrypt) Hash() error {
	salt, err := GenerateSalt(sCrypt.SaltLen)
	if err != nil {
		return err
	}

	sCrypt.Salt = salt

	dk, err := scrypt.Key([]byte(sCrypt.Plain), sCrypt.Salt, sCrypt.N, sCrypt.R, sCrypt.P, sCrypt.KeyLen)
	if err != nil {
		return err
	}
	sCrypt.DK = dk

	sCrypt.Hashed = fmt.Sprintf("%d$%d$%d$%x$%x", sCrypt.N, sCrypt.R, sCrypt.P, sCrypt.Salt, sCrypt.DK)

	return nil
}

// Validate sCrypt.Plain against sCrypt.Hashed
func (sCrypt *SCrypt) Validate() bool {
	existing, err := decodeSCryptHash(sCrypt.Hashed)
	if err != nil {
		return false
	}

	dk, err := scrypt.Key([]byte(sCrypt.Plain), existing.Salt, existing.N, existing.R, existing.P, existing.KeyLen)
	if err != nil {
		return false
	}

	if subtle.ConstantTimeCompare(existing.DK, dk) == 1 {
		return true
	}

	return false
}

// decodeSCryptHash
func decodeSCryptHash(hash string) (*SCrypt, error) {
	vals := strings.Split(hash, "$")

	// P, N, R, salt, scrypt derived key
	if len(vals) != 5 {
		return nil, errors.New("invalid hash")
	}

	sCrypt := &SCrypt{}
	var err error

	sCrypt.N, err = strconv.Atoi(vals[0])
	if err != nil {
		return nil, errors.New("invalid hash")
	}

	sCrypt.R, err = strconv.Atoi(vals[1])
	if err != nil {
		return nil, errors.New("invalid hash")
	}

	sCrypt.P, err = strconv.Atoi(vals[2])
	if err != nil {
		return nil, errors.New("invalid hash")
	}

	sCrypt.Salt, err = hex.DecodeString(vals[3])
	if err != nil {
		return nil, errors.New("invalid hash")
	}
	sCrypt.SaltLen = len(sCrypt.Salt)

	sCrypt.DK, err = hex.DecodeString(vals[4])
	if err != nil {
		return nil, errors.New("invalid hash")
	}
	sCrypt.KeyLen = len(sCrypt.DK)

	return sCrypt, nil
}
