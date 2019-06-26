package upload

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"

	"github.com/semirm-dev/go-common/storage/cache"

	"github.com/semirm-dev/go-common/api"
	"github.com/semirm-dev/go-common/jwt"
)

// Service for uploader
type Service struct {
	*Config
	api.Router
	*api.Auth
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

// Response for file uploads
type Response struct {
	Uploaded  []*Uploaded `json:"uploaded"`
	Failed    []*Failed   `json:"failed"`
	TotalSize int         `json:"-"`
	TotalTime string      `json:"-"`
}

// Initialize Service
func (service *Service) Initialize() {
	router := &api.MuxRouterAdapter{Router: mux.NewRouter()}

	if len(service.Endpoints) > 0 {
		for _, endpoint := range service.Endpoints {
			if endpoint.OnPreExecute == nil {
				endpoint.OnPreExecute = onPreExecute
			}

			if endpoint.OnFinished == nil {
				endpoint.OnFinished = onFinished
			}

			if endpoint.UseAuthMiddleware {
				auth := &api.Auth{
					Token:       service.Token,
					CacheClient: service.Cache,
				}

				router.Handle(endpoint.URL, auth.Middleware(service.uploadFunc(endpoint.Uploader, endpoint.OnPreExecute, endpoint.OnFinished)), "POST")

				auth.Apply(router)

				service.Auth = auth
			} else {
				router.HandleFunc(endpoint.URL, service.uploadFunc(endpoint.Uploader, endpoint.OnPreExecute, endpoint.OnFinished), "POST")
			}
		}
	}

	service.Router = router
}

// Listen and serve routes
func (service *Service) Listen() {
	service.Router.Listen(&api.Config{
		Host: service.Config.Host,
		Port: service.Config.Port,
	})
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

		r.ParseMultipartForm(uploader.MaxMemory)

		if r.MultipartForm == nil {
			logrus.Error("no files provided")
			return
		}

		response := &Response{
			Uploaded: make([]*Uploaded, 0),
			Failed:   make([]*Failed, 0),
		}

		var uploadProgress sync.WaitGroup

		for _, handler := range r.MultipartForm.File[uploader.MultipartForm] {
			file, err := handler.Open()
			if err != nil {
				logrus.Error("unexpected error while opening file: ", err)
				continue
			}

			uploadProgress.Add(1)
			uploaded, failed := uploader.Upload(file, handler.Filename)

			go func() {
				select {
				case u := <-uploaded:
					response.TotalSize += u.Size
					response.Uploaded = append(response.Uploaded, u)
				case f := <-failed:
					response.Failed = append(response.Failed, f)
				}

				uploadProgress.Done()
			}()
		}

		uploadProgress.Wait()

		finishTime := time.Now()
		logrus.Infof(
			"upload started at: %v | finished at: %v | finished in: %v | total size: %v",
			startTime.Format("15:04:05"),
			finishTime.Format("15:04:05.000"),
			finishTime.Sub(startTime),
			response.TotalSize,
		)

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
