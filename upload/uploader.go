package upload

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const extBytesLen = 512

// Uploader is responsible to upload files
type Uploader struct {
	Destination                string
	FileSize                   int
	AllowNonMimeTypeExtensions bool
	AllowedExtensions          []string
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
// Upload is run in separate goroutine for better performance
// Response is passed to either *Uploaded or *Failed channels, depends if the upload was successful or not
// Both *Uploaded and *Failed channels will be closed automatically on goroutine's return
//
// Destination path will be created if it doesnt exist
// FileSize will always be checked
// AllowedExtensions is optional
//
// TODO: move calls to createPathIfNotExists, fileExtension, allowedExtension to validation interface
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

		if err := uploader.fileSizeValid(fileBytes); err != nil {
			failed <- &Failed{
				File:    fileName,
				Message: err.Error(),
			}
			return
		}

		filePath, err := uploader.filePath(fileBytes, fileName)
		if err != nil {
			failed <- &Failed{
				File:    fileName,
				Message: err.Error(),
			}
			return
		}

		// performLongUploadTest(ext) // for testing only

		if err := ioutil.WriteFile(filePath, fileBytes, 0644); err != nil {
			failed <- &Failed{
				File:    fileName,
				Message: "write file failed: " + filePath,
			}
			return
		}

		finishTime := time.Now()

		uploaded <- &Uploaded{
			File: filePath,
			Size: len(fileBytes),
			Time: fmt.Sprint(finishTime.Sub(startTime)),
		}
	}()

	return uploaded, failed
}

func (uploader *Uploader) fileSizeValid(fileBytes []byte) error {
	fLen := len(fileBytes)

	if fLen > uploader.FileSize {
		return errors.New(fmt.Sprintf("file size exceeds limit by %d bytes, len[%d], maxLen[%d]", fLen-uploader.FileSize, fLen, uploader.FileSize))
	}

	if fLen < extBytesLen {
		return errors.New(fmt.Sprintf("not enough bytes to read file extension, bytes[%d]", fLen))
	}

	return nil
}

func (uploader *Uploader) filePath(fileBytes []byte, fileName string) (string, error) {
	ext, err := fileExtension(fileBytes[:extBytesLen], fileName)
	if err != nil {
		return "", errors.New("failed to get file extension: " + err.Error())
	}

	if !uploader.allowedExtension(ext) {
		return "", errors.New("file extension not allowed: " + ext)
	}

	// make sure we use valid extension from fileBytes, not user defined
	fileName = trimExtension(fileName) + ext

	return filepath.Join(uploader.Destination, fileName), err
}

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
	var extension string

	fileType := http.DetectContentType(extBytes)

	fileEndings, err := mime.ExtensionsByType(fileType)
	if err != nil {
		return "", err
	}

	// mime type extension, get extension from extBytes
	if len(fileEndings) > 0 {
		extension = fileEndings[0]
	} else { // non-mime type extension, get extension from fileName
		file := strings.Split(fileName, ".")

		extension = "." + file[len(file)-1]
	}

	return extension, nil
}

func createPathIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create upload destination directory [%v], consider creeating one manually: %v", path, err)
		}
	}

	return nil
}

func trimExtension(file string) string {
	_file := strings.Split(file, ".")

	if len(_file) < 2 {
		return file
	}

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
