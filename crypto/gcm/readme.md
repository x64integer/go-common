### GCM encryption/decryption

## ENV variables

| ENV           | Default value |
|:--------------|:-------------:|
| CRYPTO_SECRET |               |

* **Encrypt**
```
encrypted, hex, base64 := gcm.Encrypt("some input")
```

* **Decrypt**
```
decrypted, err := gcm.Decrypt("encrypted value")
```