### Usage
> NOTE: Middlewares in progress

* **Initialize router**
```
r := api.NewRouter(&api.Config{
    Host:        "localhost",
    Port:        "8080",
    WaitTimeout: time.Second * 15,
    MapRoutes: func(r api.RouteHandler) {
        r.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
            w.Write([]byte("hello :)"))
        })
    },
    // If not defined, auth routes will not be initialized
    Auth: &api.Auth{
        RegisterPath: "/register",
        LoginPath:    "/login",
        LogoutPath:   "/logout",
        Entity:       &User{},

        // Optionally, override default OnError and OnSuccess behaviour
        OnError: func(err error, w http.ResponseWriter) {
            // handle error
        },
        OnSuccess: func(payload []byte, w http.ResponseWriter) {
            // handle successful response
        },
    },
})
```

* **Start http server**
```
r.Listen()
```

```
// User entity example
type User struct {
	ID       string    `auth:"id" auth_type:"uuid"`
	Username string    `json:"user_name" auth:"username" auth_type:"credential"`
	Email    string    `auth:"email" auth_type:"credential"`
	Password string    `auth:"password" auth_type:"secret"`
	DoB      time.Time `auth:"date_of_birth"`
}
```