package crypto

import (
	"sync"
)

// Crypter for encryption and decryption
type Crypter interface {
	Encrypt([]byte) (string, string, string, error)
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
