### Usage
> NOTE: Middlewares and Auth in progress

* **Initialize router**
```
// setup storage clients
st := storage.DefaultContainer(storage.SQLClient | storage.CacheClient)

st.Connect()

gateway := &Gateway{
    Token: &jwt.Token{
        Secret: []byte(util.Env("JWT_SECRET_KEY", "some-random-string-123")),
    },
    Storage: st,
    UserAccountRepository: &my.UserAccountRepositoryImpl{
        SQL: st.SQL,
    },
}

// create router
r := api.NewRouter(&api.Config{
    Host:        "localhost",
    Port:        "8080",
    
    // optionally, setup authentication
    Auth: &api.Auth{
        Token:              gateway.Token,
        CacheClient:        gateway.Storage.Cache,
        UserAccountRepository: gateway.UserAccountRepository,
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
