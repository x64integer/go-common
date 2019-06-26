package upload

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const extBytesLen = 512

// Uploader is responsible to upload files
type Uploader struct {
	Destination                string
	FilePrefix                 string
	MultipartForm              string
	MaxMemory                  int64
	FileSize                   int
	AllowNonMimeTypeExtensions bool
	AllowedExtensions          []string
	L                          sync.Mutex
}

// Uploaded file data
type Uploaded struct {
	File string `json:"file"`
	Size int    `json:"size"`
	Time string `json:"time"`
}

// Failed upload file data
type Failed struct {
	File    string `json:"file"`
	Message string `json:"message"`
}

// Upload file
// Upload is run in separate goroutine for better performance and is thread-safe
// Response is passed to either *Uploaded or *Failed channels, depends if the upload was successful or not
// Both *Uploaded and *Failed channels will be closed automatically on goroutine's return
//
// Destination path will be created if it doesnt exist
// FileSize will always be checked
// AllowedExtensions is optional
func (uploader *Uploader) Upload(reader io.Reader, fileName string) (<-chan *Uploaded, <-chan *Failed) {
	uploaded := make(chan *Uploaded)
	failed := make(chan *Failed)

	go func() {
		defer func() {
			close(uploaded)
			close(failed)
		}()

		startTime := time.Now()

		if err := createPathIfNotExists(uploader.Destination); err != nil {
			failed <- &Failed{
				File:    fileName,
				Message: "failed to create destination path: " + uploader.Destination,
			}
			return
		}

		fileBytes, err := ioutil.ReadAll(reader)
		if err != nil {
			failed <- &Failed{
				File:    fileName,
				Message: "failed to read from file reader: " + err.Error(),
			}
			return
		}

		size := len(fileBytes) / 1024

		if len(fileBytes) > uploader.FileSize {
			failed <- &Failed{
				File:    fileName,
				Message: fmt.Sprintf("file size exceeds limit by %d bytes, len[%d], maxLen[%d]", len(fileBytes)-uploader.FileSize, len(fileBytes), uploader.FileSize),
			}
			return
		}

		if len(fileBytes) < extBytesLen {
			failed <- &Failed{
				File:    fileName,
				Message: fmt.Sprintf("not enough bytes to read file extension, bytes[%d]", len(fileBytes)),
			}
			return
		}

		ext, err := fileExtension(fileBytes[:extBytesLen], fileName)
		if err != nil {
			failed <- &Failed{
				File:    fileName,
				Message: "failed to get file extension",
			}

			return
		}

		if !uploader.allowedExtension(ext) {
			failed <- &Failed{
				File:    fileName,
				Message: "file extension not allowed: " + ext,
			}

			return
		}

		fileName = trimExtension(fileName)

		name := uploader.FilePrefix + "*-" + fileName + ext

		performLongUploadTest(ext) // for testing only

		uploadedFile, err := uploader.writeFile(fileBytes, uploader.Destination, name)
		if err != nil {
			failed <- &Failed{
				File:    fileName,
				Message: "write file failed: " + name,
			}
			return
		}

		finishTime := time.Now()

		uploaded <- &Uploaded{
			File: uploadedFile.Name(),
			Size: size,
			Time: fmt.Sprint(finishTime.Sub(startTime)),
		}
	}()

	return uploaded, failed
}

// writeFile will create file in path and write content to the file
func (uploader *Uploader) writeFile(content []byte, path string, fileName string) (*os.File, error) {
	tempFile, err := ioutil.TempFile(path, fileName)
	if err != nil {
		return nil, err
	}
	defer tempFile.Close()

	tempFile.Write(content)

	return tempFile, nil
}

// allowedExtension will check if given extension is allowed
func (uploader *Uploader) allowedExtension(extension string) bool {
	if len(uploader.AllowedExtensions) == 0 && uploader.AllowNonMimeTypeExtensions {
		return true
	}

	for _, ext := range uploader.AllowedExtensions {
		if strings.Replace(extension, ".", "", -1) == strings.Replace(ext, ".", "", -1) {
			return true
		}
	}

	return false
}

// fileExtension returns file extension with dot prefix
// .jpg, .png, .bmp, .exe, etc...
func fileExtension(extBytes []byte, fileName string) (string, error) {
	var extenstion string

	fileType := http.DetectContentType(extBytes)

	fileEndings, err := mime.ExtensionsByType(fileType)
	if err != nil {
		return "", err
	}

	// mime type extension, get extension from extBytes
	if len(fileEndings) > 0 {
		extenstion = fileEndings[0]
	} else { // non-mime type extension, get extension from fileName
		file := strings.Split(fileName, ".")

		extenstion = "." + file[len(file)-1]
	}

	return extenstion, nil
}

// createPathIfNotExists is helper function to create directory + subdirectories if such path does not exist
func createPathIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create upload destination directory [%v], consider creeating one manually: %v", path, err)
		}
	}

	return nil
}

// trimExtension from given file name
func trimExtension(file string) string {
	_file := strings.Split(file, ".")

	_file = _file[:len(_file)-1]

	return strings.Join(_file, ".")
}

func performLongUploadTest(ext string) {
	switch ext {
	case ".gif":
		time.Sleep(7 * time.Second)
	case ".jpg":
		time.Sleep(2 * time.Second)
	case ".png":
		time.Sleep(4 * time.Second)
	}
}
