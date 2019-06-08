package upload

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/semirm-dev/go-common/storage/cache"

	"github.com/semirm-dev/go-common/api"
	"github.com/semirm-dev/go-common/jwt"
)

// Service for uploader
type Service struct {
	*Config
	*api.Router
	*jwt.Token
	Cache cache.Service
	*Uploader
	OnError    func(error, http.ResponseWriter)
	OnFinished func(*Response, http.ResponseWriter)
}

// Config for router
type Config struct {
	Host          string
	Port          string
	URL           string
	UseMiddleware bool
}

// Response for file uploads
type Response struct {
	Uploaded []*Uploaded `json:"uploaded"`
	Failed   []*Failed   `json:"failed"`
}

// Initialize Service
func (service *Service) Initialize() {
	if service.OnFinished == nil {
		service.OnFinished = onFinished
	}

	if service.OnError == nil {
		service.OnError = onError
	}

	r := api.NewRouter(&api.Config{
		Host: service.Config.Host,
		Port: service.Config.Port,
	})

	if service.Config.UseMiddleware {
		r.Auth = &api.Auth{
			Token:       service.Token,
			CacheClient: service.Cache,
		}

		r.Handle(service.Config.URL, r.Auth.Middleware(service.upload), "POST")
	} else {
		r.HandleFunc(service.Config.URL, service.upload, "POST")
	}

	service.Router = r
}

// Listen and serve routes
func (service *Service) Listen() {
	service.Router.Listen()
}

// upload API endpoint will handle file upload
func (service *Service) upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(service.Uploader.FileSize)

	response := &Response{
		Uploaded: make([]*Uploaded, 0),
		Failed:   make([]*Failed, 0),
	}

	for _, handler := range r.MultipartForm.File[service.Uploader.FormFile] {
		file, err := handler.Open()
		if err != nil {
			logrus.Error("unexpected error while opening file: ", err)
			continue
		}

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			response.Failed = append(response.Failed, &Failed{
				File:    handler.Filename,
				Message: "failed to read file bytes",
			})

			continue
		}

		uploaded, err := service.Uploader.Upload(fileBytes, handler.Filename)
		file.Close()

		if err != nil {
			response.Failed = append(response.Failed, &Failed{
				File:    handler.Filename,
				Message: err.Error(),
			})

			continue
		}

		response.Uploaded = append(response.Uploaded, uploaded)
	}

	service.OnFinished(response, w)
}

// onError default callback
func onError(err error, w http.ResponseWriter) {
	logrus.Error("upload failed: ", err)
}

// onFinished default callback
func onFinished(response *Response, w http.ResponseWriter) {
	b, err := json.Marshal(response)
	if err != nil {
		logrus.Error("upload failed: ", err)
		return
	}

	w.Write(b)
}
