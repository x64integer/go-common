### Usage
> NOTE: Middlewares and Auth in progress

* **Initialize router**
```
st := storage.DefaultContainer(storage.SQLClient | storage.CacheClient)

st.Connect()

r := api.NewRouter(&api.Config{
    Host:        "localhost",
    Port:        "8080",
    
    // optionally, setup authentication
    Auth: &api.Auth{
        Token: &jwt.Token{
            Secret: []byte(util.Env("JWT_SECRET_KEY", "some-random-string-123")),
        },
        Cache: st,
        UserAccountRepository: &my.UserAccountRepositoryImpl{
            SQL: st.SQL,
        },
        PasswordResetRepository: &my.PasswordResetRepositoryImpl{
            SQL: st.SQL,
        },
    },
})

r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Spotted Gateway :)"))
}, "GET")
```

* **Start http server**
```
r.Listen()
```