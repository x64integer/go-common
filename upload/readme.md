### Usage

```
// initialize *Uploader
uploader := &upload.Uploader{
    Destination:       "./uploads",
    FileSize:          2 << 20,  // MB
    AllowedExtensions: []string{".jpg", ".png", ".bmp", ".gif"},
}

file, err := os.Open("/home/semirma/Pictures/layan-dark.png")
if err != nil {
    logrus.Info("file open failed: ", err)
}

var uploadProgress sync.WaitGroup

uploadProgress.Add(1)
// upload file with content from io.Reader
uploaded, failed := uploader.Upload(file, "myUploadedFileName.png") // .png ext will be ignored, extension is used from file bytes

// read data from uploaded and failed channels
go func() {
    select {
    case u := <-uploaded:
        logrus.Info("uploaded: ", u)
    case f := <-failed:
        logrus.Info("failed: ", f)
    }

    uploadProgress.Done()
}()

uploadProgress.Wait()
```
