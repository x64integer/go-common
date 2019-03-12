* **Initialize storage container**
> Default implementation will be good enough for most cases. However, feel free to initialize your own clients/instances
```
st := storage.DefaultContainer(storage.SQLClient | storage.RedisClient | storage.ElasticClient)
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

> Elastic example
```
entities := []*elastic.Entity{
    &elastic.Entity{
        ID: "1",
        Content: &struct {
            Name  string
            Email string
        }{
            Name:  "semir1",
            Email: "semir1@email.com",
        },
    },
    &elastic.Entity{
        ID: "2",
        Content: &struct {
            Name  string
            Email string
        }{
            Name:  "semir2",
            Email: "semir2@email.com",
        },
    },
}

if err := st.Elastic.BulkInsert(context.Background(), "users", "_doc", entities...); err != nil {
    log.Println("elasticsearch bulk insert failed: ", err)
}

entity := &elastic.Entity{
    ID: "3",
    Content: &struct {
        Name  string
        Email string
    }{
        Name:  "semir3",
        Email: "semir3@email.com",
    },
}

if err := st.Elastic.Insert(context.Background(), "users", "_doc", entity); err != nil {
    log.Println("elasticsearch insert failed: ", err)
}

search := &elastic.SearchEntity{
    Term: "semir1",
    Fields: []string{
        "Name",
    },
}

resp, err := st.Elastic.SearchByTerm(context.Background(), "users", "_doc", search)
if err != nil {
    log.Fatalln("elasticsearch failed for term: ", err)
}

log.Println(string(resp))
```