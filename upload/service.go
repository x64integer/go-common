package upload

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

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

	Endpoints []*Endpoint
}

// Config for router
type Config struct {
	Host string
	Port string
}

// Endpoint for upload router
type Endpoint struct {
	URL               string
	UseAuthMiddleware bool
	*Uploader
	OnPreExecute func(http.ResponseWriter, *http.Request) ([]byte, bool)
	OnFinished   func(*Response, http.ResponseWriter, *http.Request)
}

// Initialize Service
func (service *Service) Initialize() {
	r := api.NewRouter(&api.Config{
		Host: service.Config.Host,
		Port: service.Config.Port,
	})

	if len(service.Endpoints) > 0 {
		for _, endpoint := range service.Endpoints {
			if endpoint.OnPreExecute == nil {
				endpoint.OnPreExecute = onPreExecute
			}

			if endpoint.OnFinished == nil {
				endpoint.OnFinished = onFinished
			}

			if endpoint.UseAuthMiddleware {
				r.Auth = &api.Auth{
					Token:       service.Token,
					CacheClient: service.Cache,
				}

				r.Handle(endpoint.URL, r.Auth.Middleware(service.uploadFunc(endpoint.Uploader, endpoint.OnPreExecute, endpoint.OnFinished)), "POST")
			} else {
				r.HandleFunc(endpoint.URL, service.uploadFunc(endpoint.Uploader, endpoint.OnPreExecute, endpoint.OnFinished), "POST")
			}
		}
	}

	service.Router = r
}

// Listen and serve routes
func (service *Service) Listen() {
	service.Router.Listen()
}

// upload API endpoint will handle file upload
func (service *Service) uploadFunc(
	uploader *Uploader,
	onPreExecute func(http.ResponseWriter, *http.Request) ([]byte, bool),
	onFinished func(*Response, http.ResponseWriter, *http.Request),
) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		if b, ok := onPreExecute(w, r); !ok {
			w.Write(b)
			return
		}

		response := &Response{
			Uploaded: make([]*Uploaded, 0),
			Failed:   make([]*Failed, 0),
		}

		r.ParseMultipartForm(uploader.FileSize)

		if r.MultipartForm == nil {
			logrus.Error("no files provided")
			return
		}

		var upload sync.WaitGroup

		for _, handler := range r.MultipartForm.File[uploader.FormFile] {
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
			file.Close()

			upload.Add(1)
			uploader.Upload(fileBytes, handler.Filename, response, &upload)
		}

		upload.Wait()

		finishTime := time.Now()
		logrus.Infof(
			"upload started at: %v | finished at: %v | finished in: %v | total size: %v",
			startTime.Format("15:04:05"),
			finishTime.Format("15:04:05.000"),
			finishTime.Sub(startTime),
			response.TotalSize,
		)

		response.TotalTime = fmt.Sprint(finishTime.Sub(startTime))

		onFinished(response, w, r)
	}
}

// onFinished default callback
func onFinished(response *Response, w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(response)
	if err != nil {
		logrus.Error("upload failed: ", err)
		return
	}

	w.Write(b)
}

// onPreExecute default callback
func onPreExecute(w http.ResponseWriter, r *http.Request) ([]byte, bool) {
	return nil, true
}
