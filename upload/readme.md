### Upload service

Example

```
st := storage.DefaultContainer(storage.SQLClient | storage.CacheClient)

st.Connect()

service := &upload.Service{
    Config: &upload.Config{
        Host: util.Env("UPLOAD_SERVICE_HOST", "localhost"),
        Port: util.Env("UPLOAD_SERVICE_PORT", "9999"),
        URL:  "/upload",
        // UseMiddleware: true,
    },

    Uploader: &upload.Uploader{
        Destination: "./uploads/topics",
        FilePrefix:  "topic_", // optional
        FormFile:    "topicUpload",
        FileSize:    10 << 20, // 10MB
    },

    // required only if UseMiddleware
    Token: &jwt.Token{
        Secret: []byte(util.Env("JWT_SECRET_KEY", "some-random-string-123")),
    },
    Cache: st.Cache,
}

service.Initialize()

service.Listen()
```