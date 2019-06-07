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
    UserAccountRepo: &my.UserAccountRepositoryImpl{
        SQL: st.SQL,
    },
    PasswordResetRepo: &my.PasswordResetRepositoryImpl{
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
        RepositoryProvider: gateway,
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

* **RepositoryProvider impl**
```
// Gateway is usually wrapper for api.Router and implements api.auth.RepositoryProvider
type Gateway struct {
    *jwt.Token
    Storage *storage.Container
    UserAccountRepo   *my.UserAccountRepositoryImpl
    PasswordResetRepo *my.PasswordResetRepositoryImpl
}

// UserAccountRepository implements api.auth.RepositoryProvider.UserAccountRepository
func (gateway *Gateway) UserAccountRepository() user.Repository {
	returngateway.UserAccountRepo
}

// PasswordResetRepository implements api.auth.RepositoryProvider.PasswordResetRepository
func (gateway *Gateway) PasswordResetRepository() user.PasswordResetRepository {
	return gateway.PasswordResetRepo
}
```

Example: https://github.com/semirm-dev/spotted-gateway/blob/master/api/gateway.go
