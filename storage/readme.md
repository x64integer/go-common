* **Initialize storage container**
> Default implementation will be good enough for most cases. However, feel free to initialize your own clients/instances
```
st := storage.DefaultContainer(storage.SQLClient | storage.RedisClient | storage.ElasticClient)

st.Expose = true // so it can be accessed globally, storage.C.*, but only after st.Connect() call!
st.Cache = &custom.Cache{} // use custom cache for instance (redis by default)
```

* **Open connections and initialize clients**
```
st.Connect()
```

* **Usage**
> SQL client

| ENV          | Default value |
|:-------------|:-------------:|
| SQL_DRIVER   | postgres      |
| SQL_HOST     | localhost     |
| SQL_NAME     |               |
| SQL_PORT     | 5432          |
| SQL_USER     | postgres      |
| SQL_PASSWORD | postgres      |
| SQL_SSLMODE  | disable       |
```
rows, err := st.SQL.Query("select * from users where id = ?", 1)
if err != nil {
    log.Println("query failed: ", err)
}
defer rows.Close()
```

> Redis

| ENV            | Default value |
|:---------------|:-------------:|
| REDIS_HOST     |               |
| REDIS_PORT     | 6379          |
| REDIS_DB       | 0             |
| REDIS_PASSWORD |               |
```
if err := st.Redis.Store(&redis.Item{
    Key:        "semir_redis",
    Value:      util.RandomStr(16),
    Expiration: time.Hour,
}); err != nil {
    log.Println("redis failed to store: ", err)
}
```

> Cache (redis by default)
```
if err := st.Cache.Store(&cache.Item{
    Key:        "semir_cache",
    Value:      util.RandomStr(16),
    Expiration: time.Hour,
}); err != nil {
    log.Println("cache failed to store: ", err)
}
```

> Elasticsearch

| ENV          | Default value |
|:-------------|:-------------:|
| ELASTIC_HOST | 127.0.0.1     |
| ELASTIC_PORT | 9200          |
```
// BULK INSERT
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

// INSERT
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

// SEARCH BY TERM
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

> Cassandra

| ENV                | Example                                      |
|:-------------------|:--------------------------------------------:|
| CASSANDRA_KEYSPACE | some_keyspace_name                           |
| CASSANDRA_HOSTS    | 127.0.0.1,127.0.0.2,127.0.0.3                |
| CASSANDRA_USERNAME | set only if it's required from cluster setup |
| CASSANDRA_PASSWORD | set only if it's required from cluster setup |
```
// INSERT
if err := st.Cassandra.Exec(
    "INSERT INTO users (id, name, email, posted_time) VALUES (?, ?, ?, ?)", gocql.TimeUUID(), "semir", "semir@email.com", time.Now(),
); err != nil {
    log.Println("cassandra insert failed: ", err)
}

// UPDATE
if err := st.Cassandra.Exec(
    "UPDATE users SET email = ?, posted_time = ? WHERE id = ?", "semir@email.com_updated", time.Now(), "9a0b645e-44c2-11e9-8cce-acde48001122",
); err != nil {
    log.Println("cassandra update failed: ", err)
}

// DELETE
if err := st.Cassandra.Exec(
    "DELETE FROM users WHERE id = ?", "2482370a-44c5-11e9-b9bb-acde48001122",
); err != nil {
    log.Println("cassandra delete failed: ", err)
}

// SELECT
iterator := st.Cassandra.Select(
    "SELECT name, email FROM users WHERE id = ?", "9a0b645e-44c2-11e9-8cce-acde48001122",
)

m := map[string]interface{}{}
for iterator.MapScan(m) {
    log.Println(m)
}

var name, email string
for iterator.Scan(&name, &email) {
    log.Println(name, email)
}
```
