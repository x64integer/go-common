### Password hashing and validation

* **Hash password**
```
hashed, err := password.Hash("my password")
```

* **Validate hashed password**
```
valid := password.Valid("hashed password", "input clean password")
```