package crypto

import (
	"errors"
	"sync"
)

// ErrMissingCrypter error
var ErrMissingCrypter = errors.New("missing Crypter implementation")

// Crypter for encryption and decryption
type Crypter interface {
	// Encrypt returns original encrypted value, hex, base64 and error
	Encrypt([]byte) (string, string, string, error)
	// Decrypt returns original message and error
	Decrypt(string) (string, error)
}

// Cipher ...
type Cipher struct {
	Crypter
	Encrypted   string
	Decrypted   string
	Hex         string
	Base64      string
	EncryptLock sync.Mutex
	DecryptLock sync.Mutex
}

// Encrypt given payload
func (cipher *Cipher) Encrypt(payload []byte) error {
	cipher.EncryptLock.Lock()
	defer cipher.EncryptLock.Unlock()

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
	cipher.DecryptLock.Lock()
	defer cipher.DecryptLock.Unlock()

	dec, err := cipher.Crypter.Decrypt(encrypted)
	if err != nil {
		return nil
	}

	cipher.Decrypted = dec

	return nil
}
