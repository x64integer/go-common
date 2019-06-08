package upload

import (
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"strings"
)

// Uploader is responsible to upload files
type Uploader struct {
	Destination              string
	FilePrefix               string
	FormFile                 string
	FileSize                 int64
	AllowExtensionExceptions bool
}

// Uploaded contains uploaded file infomration
type Uploaded struct {
	File string `json:"file"`
}

// Upload file
func (uploader *Uploader) Upload(fileBytes []byte, file string) (*Uploaded, error) {
	if err := createPathIfNotExists(uploader.Destination); err != nil {
		return nil, err
	}

	fileExtension, err := uploader.fileExtension(fileBytes)
	if err != nil {
		return nil, err
	}

	if fileExtension != "" {
		// in case of AllowExtensionExceptions (executables and non-media types)
		file = trimExtension(file)
	}

	fileName := uploader.FilePrefix + "*-" + file + fileExtension

	uploadedFile, err := uploader.writeFile(fileBytes, uploader.Destination, fileName)
	if err != nil {
		return nil, err
	}

	return &Uploaded{
		File: uploadedFile.Name(),
	}, nil
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

// fileExtension returns .jpg, .png, etc...
func (uploader *Uploader) fileExtension(fileBytes []byte) (string, error) {
	fileType := http.DetectContentType(fileBytes)

	fileEndings, err := mime.ExtensionsByType(fileType)
	if err != nil {
		return "", err
	}

	if len(fileEndings) > 0 {
		return fileEndings[0], nil
	}

	if uploader.AllowExtensionExceptions && len(fileEndings) < 1 {
		return "", nil
	}

	return "", errors.New("invalid or failed to get file extension")
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

func trimExtension(file string) string {
	_file := strings.Split(file, ".")

	_file = _file[:len(_file)-1]

	return strings.Join(_file, "")
}
