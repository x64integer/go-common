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
	Destination                string
	FilePrefix                 string
	FormFile                   string
	FileSize                   int64
	AllowNonMimeTypeExtensions bool
	AllowedExtensions          []string
}

// Uploaded contains uploaded file infomration
type Uploaded struct {
	File string `json:"file"`
}

// Failed upload
type Failed struct {
	File    string `json:"file"`
	Message string `json:"message"`
}

// Upload file
func (uploader *Uploader) Upload(fileBytes []byte, file string) (*Uploaded, error) {
	if err := createPathIfNotExists(uploader.Destination); err != nil {
		return nil, err
	}

	fileExtension, err := uploader.fileExtension(fileBytes, file)
	if err != nil {
		return nil, err
	}

	if !uploader.allowedExtension(fileExtension) {
		return nil, errors.New("file extension not allowed: " + fileExtension)
	}

	file = trimExtension(file)

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
