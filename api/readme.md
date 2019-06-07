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
        Authenticator: &Gateway{
            Storage: st,
            Token: &jwt.Token{
                Secret: []byte(util.Env("JWT_SECRET_KEY", "some-random-string-123")),
            },
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

* **Authenticator impl**
```
// Gateway is usually wrapper for api.Router and implements api.auth.Authenticator
type Gateway struct {
	Storage *storage.Container
	*jwt.Token
}

// UserAccountRepository implements api.auth.Authenticator.UserAccountRepository
func (gateway *Gateway) UserAccountRepository() user.Repository {
	return &my.UserAccountRepositoryImpl{
		SQL: gateway.Storage.SQL,
	}
}

// PasswordResetRepository implements api.auth.Authenticator.PasswordResetRepository
func (gateway *Gateway) PasswordResetRepository() user.PasswordResetRepository {
	return &my.PasswordResetRepositoryImpl{
		SQL: gateway.Storage.SQL,
	}
}

// JWT implements api.auth.Authenticator.JWT
func (gateway *Gateway) JWT() *jwt.Token {
	return gateway.Token
}

// CacheClient implements api.auth.Authenticator.CacheClient
func (gateway *Gateway) CacheClient() cache.Service {
	return gateway.Storage.Cache
}
```

Example: https://github.com/semirm-dev/spotted-gateway/blob/master/api/gateway.go