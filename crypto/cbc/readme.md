### CBC encryption/decryption

## ENV variables

| ENV           | Default value |
|:--------------|:-------------:|
| CRYPTO_SECRET |               |

* **Encrypt**
```
encrypted, hex, base64 := cbc.Encrypt("some input")
```

* **Decrypt**
```
decrypted, err := cbc.Decrypt("encrypted value")
```