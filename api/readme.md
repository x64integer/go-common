### Usage
> NOTE: Prototype, in progress

```
// setup storage cache client, needed for Auth
st := storage.DefaultContainer(storage.CacheClient)

st.Connect()

// create router
router := &_api.MuxRouterAdapter{Router: mux.NewRouter()}
// router := &_api.IrisRouterAdapter{Application: iris.Default()}

auth := &_api.Auth{
    Token: &jwt.Token{
        Secret: []byte("my-random-string-123"),
    },
    CacheClient:             st.Cache,
    UserAccountRepository:   &my.UserAccountRepositoryImpl{},
    PasswordResetRepository: &my.PasswordResetRepositoryImpl{},
    ServiceURL:              "localhost:8080",
    RequireConfirmation:     true,
}

auth.Apply(router)

// define routes
router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Home :)"))
}, "GET")

router.Listen(&_api.Config{
    Host: "localhost",
    Port: "8080",
})
```

Example: https://github.com/semirm-dev/spotted-gateway/blob/master/api/gateway.go
