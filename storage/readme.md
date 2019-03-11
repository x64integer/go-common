* **Initialize storage container**
> Default implementation will be good enough for most cases. However, feel free to initialize your own clients/instances
```
st := storage.DefaultContainer(storage.SQLClient | storage.RedisClient)
st.Cache = &custom.Cache{} // use custom cache for instance (redis by default)
```

* **Open connections and initialize clients**
```
st.Connect()
```

* **Usage**
> SQL client
```
rows, err := st.SQL.Query("select * from users where id = ?", 1)
if err != nil {
    log.Println("query failed: ", err)
}
defer rows.Close()
```

> Redis example
```
if err := st.Redis.Store(&redis.Item{
    Key:        "semir_redis",
    Value:      util.RandomStr(16),
    Expiration: time.Hour,
}); err != nil {
    log.Println("redis failed to store: ", err)
}
```

> Cache example (redis by default)
```
if err := st.Cache.Store(&cache.Item{
    Key:        "semir_cache",
    Value:      util.RandomStr(16),
    Expiration: time.Hour,
}); err != nil {
    log.Println("cache failed to store: ", err)
}
```