### JWT

* **Generate jwt with claims**
```
token := &jwt.Token{
    Secret: []byte("testkey"),
}

if err := token.Generate(&jwt.Claims{
    Expiration: time.Hour * 24,
    Fields: map[string]interface{}{
        "username": "semir",
        "email":    "semir@mail.com",
        "id":       "semir-123",
    },
}); err != nil {
    log.Fatalln("failed to generate jwt: ", err)
}
```

* **Validate and get jwt claims**
```
claims, valid := token.ValidateAndExtract(token.Content)
```