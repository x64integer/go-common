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

    // OnFinished: func(response *upload.Response, w http.ResponseWriter) {
    // 	b, err := json.Marshal(response.Uploaded)
    // 	if err != nil {
    // 		logrus.Error("failed to marshal response.Uploaded: ", err)
    // 		return
    // 	}

    // 	w.Write(b)
    // },

    Uploader: &upload.Uploader{
        Destination: "./uploads/topics",
        FilePrefix:  "topic_", // optional
        FormFile:    "topicUpload",
        FileSize:    32 << 20, // 32MB
        // AllowExtensionExceptions: true,
        AllowedExtensions:          []string{".jpg", ".png", ".exe"},
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