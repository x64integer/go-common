### Usage
> NOTE: Middlewares in progress

* **Initialize router and start listening**
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
        LoginPath:    "/login", // this is can be removed since it's overriden in custom Login.Path
        // Use &User{} entity for both register and login
        Entity: &User{},
        // Optionally, use different entities for register and login
        // Not changed values will default to Entity: &User{}
        Registration: &api.Registration{
            OnSuccess: func(payload []byte, w http.ResponseWriter) {
                log.Println("registration successful")
            },
        },
        Login: &api.Login{
            Path:   "/user/login", // override previous value -> /login
            Entity: &LoginUser{},
            OnError: func(err error, w http.ResponseWriter) {
                log.Println("login error occured")
                w.Write([]byte(err.Error()))
            },
            OnSuccess: func(payload []byte, w http.ResponseWriter) {
                log.Println("login successful")
                w.Write(payload)
            },
        },
    },
})

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