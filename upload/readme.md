### Usage

```
// initialize *Uploader
uploader := &upload.Uploader{
    Destination:       "./uploads/topic",
    FileSize:          2 << 20,  // MB
    AllowedExtensions: []string{".jpg", ".png", ".bmp", ".gif"},
}

var uploadProgress sync.WaitGroup

uploadProgress.Add(1)
// upload file with content from io.Reader
uploaded, failed := uploader.Upload(myReader, "myFileName.txt")

// read data from uploaded and failed channels
go func() {
    select {
    case u := <-uploaded:
        // handle uploaded file
    case f := <-failed:
        // handle failed upload
    }

    uploadProgress.Done()
}()

uploadProgress.Wait()

```
