### Password hashing and validation

> NOTE: Soon to be replaced with password.SCrypt and password.Argon2

* **Hash password**
```
hashed, err := password.Hash("my password")
```

* **Validate hashed password**
```
valid := password.Valid("hashed password", "input clean password")
```
