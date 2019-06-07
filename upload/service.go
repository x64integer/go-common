package upload

import (
	"encoding/json"
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
	OnError   func(error, http.ResponseWriter)
	OnSuccess func([]byte, http.ResponseWriter)
}

// Config for router
type Config struct {
	Host          string
	Port          string
	URL           string
	UseMiddleware bool
}

// Initialize Service
func (service *Service) Initialize() {
	if service.OnSuccess == nil {
		service.OnSuccess = onSuccess
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

	file, handler, err := r.FormFile(service.Uploader.FormFile)
	if err != nil {
		service.OnError(err, w)
		return
	}
	defer file.Close()

	uploaded, err := service.Uploader.Upload(file, handler.Filename)
	if err != nil {
		service.OnError(err, w)
		return
	}

	b, err := json.Marshal(uploaded)
	if err != nil {
		service.OnError(err, w)
		return
	}

	service.OnSuccess(b, w)
}

// onError default callback
func onError(err error, w http.ResponseWriter) {
	logrus.Error("upload failed: ", err)
	w.Write([]byte("upload failed"))
}

// onSuccess default callback
func onSuccess(payload []byte, w http.ResponseWriter) {
	w.Write(payload)
}
