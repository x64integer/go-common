### JWT

## ENV variables

| ENV       | Default value      |
|:----------|:------------------:|
| JWT_LOGIN | util.RandomStr(64) |

* **Generate jwt with claims**
```
token, err := jwt.Generate(map[string]string{
    "id":       "semir-123",
    "username": "semir",
    "email":    "semir@mail.com",
})
```

* **Validate and get jwt claims**
```
claims, valid := jwt.ValidateAndExtract(token)
```