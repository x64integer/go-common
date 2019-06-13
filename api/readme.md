### Usage
> NOTE: Auth in progress still

* **Initialize router**
```
// setup storage cache client, needed for Auth
st := storage.DefaultContainer(storage.CacheClient)

st.Connect()

// create router
r := api.NewRouter(&api.Config{
    Host:        "localhost",
    Port:        "8080",
    
    // optionally, setup authentication
    Auth: &api.Auth{
        Token: &jwt.Token{
            Secret: []byte(util.Env("JWT_SECRET_KEY", "some-random-string-123")),
        },
        CacheClient:        st.Cache,
        UserAccountRepository: &my.UserAccountRepositoryImpl{},
        RequireConfirmation:   true,
    },
})

// define routes
r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Spotted Gateway :)"))
}, "GET")
```

* **Start http server**
```
r.Listen()
```

Example: https://github.com/semirm-dev/spotted-gateway/blob/master/api/gateway.go
