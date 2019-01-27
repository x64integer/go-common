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
})

r.Listen()
```