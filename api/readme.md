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
        // 1.) Simple register, login, logout configurations
        RegisterPath: "/register",
        LoginPath:    "/login",
        LogoutPath:   "/logout",
        Entity:       &User{}, // use &User{} entity for register, login, logout

        // 2.) Custom register, login, logout configurations
        // If both simple and custom configurations are used, custom configurations have higher priority
        Registration: &api.Registration{
            Entity: &RegisterUser{}, // override auth.Entity for Registration
        },
        Login: &api.Login{
            Path:   "/user/login", // override auth.LoginPath value
            Entity: &LoginUser{},  // override auth.Entity for Login
            OnError: func(err error, w http.ResponseWriter) { // override default OnError()
                log.Println("login error occured")
                w.Write([]byte(err.Error()))
            },
            OnSuccess: func(payload []byte, w http.ResponseWriter) { // override default OnSuccess()
                log.Println("login successful")
                w.Write(payload)
            },
        },
        Logout: &api.Logout{
            Entity: &LogoutUser{}, // override auth.Entity for Logout
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