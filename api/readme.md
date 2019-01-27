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
        // common/defaults for both register and login (can be overriden by customizing api.Registration and api.Login)
        RegisterPath: "/register",
        LoginPath:    "/login", // overridden in custom Login.Path
        // Use &User{} as default entity for both register and login
        Entity: &User{}, // will not be used because it's overriden in both api.Registration and api.Login

        // Optionally, we can use different entities for register and login with customizations
        // Not changed values will be ignored
        Registration: &api.Registration{
            Entity: &RegisterUser{}, // override default entity &User{}
        },
        Login: &api.Login{
            Path:   "/user/login", // override previous value -> /login
            Entity: &LoginUser{},  // override default entity &User{}
            OnError: func(err error, w http.ResponseWriter) { // override default OnError()
                log.Println("login error occured")
                w.Write([]byte(err.Error()))
            },
            OnSuccess: func(payload []byte, w http.ResponseWriter) { // override default OnSuccess()
                log.Println("login successful")
                w.Write(payload)
            },
        },
    },
})
```

* **Start http server**
```
r.Listen()
```

* **Entities examples**
```
// User entity example
type User struct {
	ID       string    `auth:"id" auth_type:"uuid"`
	Username string    `auth:"username" auth_type:"credential"`
	Email    string    `auth:"email" auth_type:"credential"`
	Password string    `auth:"password" auth_type:"secret"`
	DoB      time.Time `auth:"date_of_birth"`
}

// LoginUser entity example
type LoginUser struct {
	Email    string `auth:"email" auth_type:"credential"`
	Password string `auth:"password" auth_type:"secret"`
}
```