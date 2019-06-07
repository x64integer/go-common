package upload

import (
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Uploader is responsible to upload files
type Uploader struct {
	Destination string
	FilePrefix  string
	FormFile    string
	FileSize    int64
}

// Uploaded after successful file upload
type Uploaded struct {
	File string `json:"file"`
}

// Upload file
func (uploader *Uploader) Upload(fileReader io.Reader, fileName string) (*Uploaded, error) {
	if _, err := os.Stat(uploader.Destination); os.IsNotExist(err) {
		if err := os.MkdirAll(uploader.Destination, os.ModePerm); err != nil {
			logrus.Errorf(
				"failed to create upload destination directory [%v], consider creeating one manually",
				uploader.Destination,
			)

			return nil, err
		}
	}

	fileBytes, err := ioutil.ReadAll(fileReader)
	if err != nil {
		return nil, err
	}

	fileExtension, err := uploader.fileExtension(fileBytes)

	file := uploader.FilePrefix + "*-" + strings.TrimSuffix(fileName, fileExtension) + fileExtension

	tempFile, err := ioutil.TempFile(uploader.Destination, file)
	if err != nil {
		return nil, err
	}
	defer tempFile.Close()

	tempFile.Write(fileBytes)

	return &Uploaded{
		File: tempFile.Name(),
	}, nil
}

// fileExtension returns .jpg, .png, etc...
func (uploader *Uploader) fileExtension(fileBytes []byte) (string, error) {
	fileType := http.DetectContentType(fileBytes)

	fileEndings, err := mime.ExtensionsByType(fileType)
	if err != nil {
		return "", err
	}

	return fileEndings[0], nil
}
