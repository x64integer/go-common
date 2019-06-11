package upload

import (
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// Uploader is responsible to upload files
type Uploader struct {
	Destination                string
	FilePrefix                 string
	FormFile                   string
	FileSize                   int64
	AllowNonMimeTypeExtensions bool
	AllowedExtensions          []string
}

// Response for file uploads
type Response struct {
	TotalSize float32     `json:"total_size"`
	TotalTime string      `json:"total_time"`
	Uploaded  []*Uploaded `json:"uploaded"`
	Failed    []*Failed   `json:"failed"`
}

// Uploaded contains uploaded file infomration
type Uploaded struct {
	File string `json:"file"`
	Size string `json:"size"`
	Time string `json:"time"`
}

// Failed upload
type Failed struct {
	File    string `json:"file"`
	Size    string `json:"size"`
	Message string `json:"message"`
}

// Upload file
//
// TODO: get rid of *Response parameter, return chan *Response instead
func (uploader *Uploader) Upload(fileBytes []byte, file string, response *Response, upload *sync.WaitGroup) {
	go func() {
		defer upload.Done()

		startTime := time.Now()

		size := float32(len(fileBytes)) / 1024

		if err := createPathIfNotExists(uploader.Destination); err != nil {
			response.Failed = append(response.Failed, &Failed{
				File:    file,
				Message: err.Error(),
				Size:    fmt.Sprint(size) + " KB",
			})

			return
		}

		fileExtension, err := uploader.fileExtension(fileBytes, file)
		if err != nil {
			response.Failed = append(response.Failed, &Failed{
				File:    file,
				Message: err.Error(),
				Size:    fmt.Sprint(size) + " KB",
			})

			return
		}

		if !uploader.allowedExtension(fileExtension) {
			response.Failed = append(response.Failed, &Failed{
				File:    file,
				Message: "file extension not allowed: " + fileExtension,
				Size:    fmt.Sprint(size) + " KB",
			})

			return
		}

		file = trimExtension(file)

		fileName := uploader.FilePrefix + "*-" + file + fileExtension

		uploadedFile, err := uploader.writeFile(fileBytes, uploader.Destination, fileName)
		if err != nil {
			response.Failed = append(response.Failed, &Failed{
				File:    fileName,
				Message: err.Error(),
				Size:    fmt.Sprint(size) + " KB",
			})

			return
		}

		finishTime := time.Now()

		response.Uploaded = append(response.Uploaded, &Uploaded{
			File: uploadedFile.Name(),
			Size: fmt.Sprint(size) + " KB",
			Time: fmt.Sprint(finishTime.Sub(startTime)),
		})

		response.TotalSize += size
	}()
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

// fileExtension returns file extension with dot prefix
// .jpg, .png, .bmp, .exe, etc...
func (uploader *Uploader) fileExtension(fileBytes []byte, fileName string) (string, error) {
	var extenstion string

	fileType := http.DetectContentType(fileBytes)

	fileEndings, err := mime.ExtensionsByType(fileType)
	if err != nil {
		return "", err
	}

	// mime type extension, get extension from fileBytes
	if len(fileEndings) > 0 {
		extenstion = fileEndings[0]
	} else { // non-mime type extension, get extension from fileName
		file := strings.Split(fileName, ".")

		extenstion = "." + file[len(file)-1]
	}

	return extenstion, nil
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
