### Password hashing and validation

* **Argon2**
```
argon := password.NewArgon2()
argon.Plain = "pwd-123"

if err := argon.Hash(); err != nil {
    logrus.Fatal("argon failed to hash pwd: ", err)
}

if argon.Validate() {
    // valid password
}
```

* **SCrypt**
```
sCrypt := password.NewSCrypt()
sCrypt.Plain = "pwd-123"

if err := sCrypt.Hash(); err != nil {
    logrus.Fatal("scrypt failed to hash pwd: ", err)
}

if sCrypt.Validate() {
    // valid password
}
```

* **BCrypt**
```
bCrypt := password.NewBCrypt()
bCrypt.Plain = "pwd-123"

if err := bCrypt.Hash(); err != nil {
    logrus.Fatal("bCrypt failed to hash pwd: ", err)
}

if bCrypt.Validate() {
    // valid password
}
```