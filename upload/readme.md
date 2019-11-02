### Usage

```
// initialize *Uploader
uploader := &upload.Uploader{
    Destination:       "./uploads/topic",
    FileSize:          2 << 20,  // MB
    AllowedExtensions: []string{".jpg", ".png", ".bmp", ".gif"},
}

response := &upload.Response{
    Uploaded: make([]*upload.Uploaded, 0),
    Failed:   make([]*upload.Failed, 0),
}

var uploadProgress sync.WaitGroup

uploadProgress.Add(1)
// upload file with content from io.Reader
uploaded, failed := uploader.Upload(myReader, "myFileName.txt")

// read data from uploaded and failed channels
go func() {
    select {
    case u := <-uploaded:
        response.Uploaded = append(response.Uploaded, u)
    case f := <-failed:
        response.Failed = append(response.Failed, f)
    }

    uploadProgress.Done()
}()

uploadProgress.Wait()

// do something with *Response
```
