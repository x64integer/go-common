package api

import (
	"net/http"
	"strings"

	"github.com/kataras/iris"
)

// IrisRouterAdapter implements Router interface
type IrisRouterAdapter struct {
	*iris.Application
}

// Listen implements Router.Listen()
func (adapter *IrisRouterAdapter) Listen(config *Config) {
	adapter.Run(iris.Addr(config.Host + ":" + config.Port))
}

// Handle implements Router.Handle()
func (adapter *IrisRouterAdapter) Handle(path string, handler http.Handler, methods ...string) {
	for _, method := range methods {
		switch strings.ToLower(method) {
		case "get":
			adapter.Get(path, func(ctx iris.Context) {
				handler.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
			})
		case "post":
			adapter.Post(path, func(ctx iris.Context) {
				handler.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
			})
		case "put":
			adapter.Put(path, func(ctx iris.Context) {
				handler.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
			})
		case "patch":
			adapter.Patch(path, func(ctx iris.Context) {
				handler.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
			})
		case "options":
			adapter.Options(path, func(ctx iris.Context) {
				handler.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
			})
		default:
			adapter.Get(path, func(ctx iris.Context) {
				handler.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
			})
		}
	}
}

// HandleFunc implements Router.HandleFunc()
func (adapter *IrisRouterAdapter) HandleFunc(path string, f func(http.ResponseWriter, *http.Request), methods ...string) {
	for _, method := range methods {
		switch strings.ToLower(method) {
		case "get":
			adapter.Get(path, func(ctx iris.Context) {
				f(ctx.ResponseWriter(), ctx.Request())
			})
		case "post":
			adapter.Post(path, func(ctx iris.Context) {
				f(ctx.ResponseWriter(), ctx.Request())
			})
		case "put":
			adapter.Put(path, func(ctx iris.Context) {
				f(ctx.ResponseWriter(), ctx.Request())
			})
		case "patch":
			adapter.Patch(path, func(ctx iris.Context) {
				f(ctx.ResponseWriter(), ctx.Request())
			})
		case "options":
			adapter.Options(path, func(ctx iris.Context) {
				f(ctx.ResponseWriter(), ctx.Request())
			})
		default:
			adapter.Get(path, func(ctx iris.Context) {
				f(ctx.ResponseWriter(), ctx.Request())
			})
		}
	}
}

// Vars implements Router.Vars()
func (adapter *IrisRouterAdapter) Vars(r *http.Request) map[string]string {
	return nil
}
