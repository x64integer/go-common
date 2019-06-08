### Upload service

Example

```
st := storage.DefaultContainer(storage.SQLClient | storage.CacheClient)

st.Connect()

service := &upload.Service{
    Config: &upload.Config{
        Host: util.Env("UPLOAD_SERVICE_HOST", "localhost"),
        Port: util.Env("UPLOAD_SERVICE_PORT", "9999"),
    },

    Token: &jwt.Token{
        Secret: []byte(util.Env("API_JWT_SECRET_KEY", "some-random-backup-string-123")),
    },
    Cache: st.Cache,
}

topicUploadEndpoint := &upload.Endpoint{
    URL:           "/upload/topic",
    UseMiddleware: true,
    Uploader: &upload.Uploader{
        Destination:       "./uploads/topic",
        FilePrefix:        "topic_",
        FormFile:          "topicUpload",
        FileSize:          32 << 20, // MB
        AllowedExtensions: []string{".jpg", ".png", ".bmp", ".gif"},
    },
    OnPreExecute: func(w http.ResponseWriter, r *http.Request) ([]byte, bool) {
        _, err := strconv.Atoi(r.Header.Get("topic_id"))
        if err != nil {
            return []byte("failed to get topic_id from headers: " + err.Error()), false
        }

        return nil, true
    },
    OnFinished: func(response *upload.Response, w http.ResponseWriter, r *http.Request) {
        topicID, err := strconv.Atoi(r.Header.Get("topic_id"))
        if err != nil {
            w.Write([]byte(fmt.Sprint("failed to get topic_id from headers: ", err)))
            return
        }

        userID, _, _, err := service.Auth.Extract(r)
        if err != nil {
            w.Write([]byte(fmt.Sprint("failed to get userID from auth: ", err)))
            return
        }

        // do something with topicID, userID, response data

        b, err := json.Marshal(response)
        if err != nil {
            logrus.Error("failed to marshal response.Uploaded: ", err)
            return
        }

        w.Write(b)
    },
}

profileUploadEndpoint := &upload.Endpoint{
    URL:           "/upload/profile",
    UseMiddleware: true,
    Uploader: &upload.Uploader{
        Destination:       "./uploads/profile",
        FilePrefix:        "profile_",
        FormFile:          "profileUpload",
        FileSize:          10 << 20, // MB
        AllowedExtensions: []string{".jpg", ".png"},
    },
}

service.Endpoints = append(service.Endpoints, topicUploadEndpoint, profileUploadEndpoint)

service.Initialize()

service.Listen()
```