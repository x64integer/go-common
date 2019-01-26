### Usage

* **Initialize router and start listening**
```
r := api.NewRouter(&api.Config{
    Host:        "localhost",
    Port:        "8080",
    WaitTimeout: time.Second * 15,
})

r.Listen()
```