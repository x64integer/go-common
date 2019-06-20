### Usage
> NOTE: Auth in progress still

* **Initialize router**
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
    // RequireConfirmation:   true,
}

auth.ServiceURL = gateway.Config.Host + ":" + gateway.Config.Port
auth.Apply(router)

// define routes
router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Spotted Gateway :)"))
}, "GET")

router.Listen(&_api.Config{
    Host: gateway.Config.Host,
    Port: gateway.Config.Port,
})
```

Example: https://github.com/semirm-dev/spotted-gateway/blob/master/api/gateway.go
