### Usage
> NOTE: Auth in progress still

```
// setup storage cache client, needed for Auth
st := storage.DefaultContainer(storage.CacheClient)

st.Connect()

// create router
router := &_api.MuxRouterAdapter{Router: mux.NewRouter()}
// router := &_api.IrisRouterAdapter{Application: iris.Default()}

auth := &_api.Auth{
    Token:                 gateway.Token,
    CacheClient:           gateway.Storage.Cache,
    UserAccountRepository: gateway.UserAccountRepository,
    ServiceURL:            "localhost:8080",
    RequireConfirmation:   true,
}

auth.Apply(router)

// define routes
router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Spotted Gateway :)"))
}, "GET")

router.Listen(&_api.Config{
    Host: "localhost",
    Port: "8080",
})
```

Example: https://github.com/semirm-dev/spotted-gateway/blob/master/api/gateway.go
