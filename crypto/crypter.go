package crypto

import (
	"errors"
)

// ErrMissingCrypter error
var ErrMissingCrypter = errors.New("missing Crypter implementation")

// ErrMissingPlain error
var ErrMissingPlain = errors.New("missing Plain property")

// Crypter for encryption and decryption
type Crypter interface {
	// Encrypt returns original encrypted value, hex, base64 and error
	Encrypt(payload []byte) (string, string, string, error)
	// Decrypt returns original message and error
	Decrypt(encrypted string) (string, error)
}

// Cipher ...
type Cipher struct {
	Crypter
	Encrypted   string
	Decrypted   string
	Hex         string
	Base64      string
}

// Encrypt given payload
func (cipher *Cipher) Encrypt(payload []byte) error {
	if cipher.Crypter == nil {
		return ErrMissingCrypter
	}

	enc, hex, b64, err := cipher.Crypter.Encrypt(payload)
	if err != nil {
		return err
	}

	cipher.Encrypted = enc
	cipher.Hex = hex
	cipher.Base64 = b64

	return nil
}

// Decrypt encrypted payload
func (cipher *Cipher) Decrypt(encrypted string) error {
	dec, err := cipher.Crypter.Decrypt(encrypted)
	if err != nil {
		return nil
	}

	cipher.Decrypted = dec

	return nil
}
