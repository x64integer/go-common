### Encryption

* **GCM encryption mode**
```
// Secret key must be compatible with AES-128 or AES-256 (or AES-512)
cipher := &crypto.Cipher{
    Crypter: &crypto.GCM{
        Secret: "test-key-1234567",
    },
}

if err := cipher.Encrypt([]byte("test")); err != nil {
    log.Fatalln("failed to encrypt: ", err)
}

log.Printf("\nEncrypted: %s\nDecrypted: %s\nHex: %s\nBase64: %s\n", cipher.Encrypted, cipher.Decrypted, cipher.Hex, cipher.Base64)

if err := cipher.Decrypt(cipher.Encrypted); err != nil {
    log.Fatalln("failed to decrypt: ", err)
}

log.Printf("\nEncrypted: %s\nDecrypted: %s\nHex: %s\nBase64: %s\n", cipher.Encrypted, cipher.Decrypted, cipher.Hex, cipher.Base64)
```

* **CBC encryption mode**
```
// Secret key and initialization vector must be compatible with AES-128 or AES-256 (or AES-512)
cipher := &crypto.Cipher{
    Crypter: &crypto.CBC{
        Secret: "test-key-1234567",
        IV:     "test-iv-12345678",
    },
}

if err := cipher.Encrypt([]byte("test")); err != nil {
    log.Fatalln("failed to encrypt: ", err)
}

log.Printf("\nEncrypted: %s\nDecrypted: %s\nHex: %s\nBase64: %s\n", cipher.Encrypted, cipher.Decrypted, cipher.Hex, cipher.Base64)

if err := cipher.Decrypt(cipher.Encrypted); err != nil {
    log.Fatalln("failed to decrypt: ", err)
}

log.Printf("\nEncrypted: %s\nDecrypted: %s\nHex: %s\nBase64: %s\n", cipher.Encrypted, cipher.Decrypted, cipher.Hex, cipher.Base64)
```

### Hashing

* **Argon2**
```
argon := crypto.NewArgon2()
argon.Plain = "value to hash"

if err := argon.Hash(); err != nil {
    logrus.Fatal("argon failed to hash Plain: ", err)
}

if argon.Validate() {
    // hash valid
}
```

* **SCrypt**
```
sCrypt := crypto.NewSCrypt()
sCrypt.Plain = "value to hash"

if err := sCrypt.Hash(); err != nil {
    logrus.Fatal("scrypt failed to hash Plain: ", err)
}

if sCrypt.Validate() {
    // hash valid
}
```

* **BCrypt**
```
bCrypt := crypto.NewBCrypt()
bCrypt.Plain = "value to hash"

if err := bCrypt.Hash(); err != nil {
    logrus.Fatal("bCrypt failed to hash Plain: ", err)
}

if bCrypt.Validate() {
    // hash valid
}
```