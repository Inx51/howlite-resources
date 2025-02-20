package api

import (
	"net/http"
	"strconv"

	"github.com/inx51/howlite/resources/api/handler"
	"github.com/inx51/howlite/resources/resource/repository"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func SetupHandlers(repository *repository.Repository) {
	http.DefaultServeMux = http.NewServeMux()

	http.Handle("GET /", otelhttp.NewHandler(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		handler.GetResource(resp, req, repository)
	}), "GetResource"))

	http.Handle("POST /", otelhttp.NewHandler(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		handler.CreateResource(resp, req, repository)
	}), "CreateResource"))

	http.Handle("HEAD /", otelhttp.NewHandler(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		handler.ResourceExists(resp, req, repository)
	}), "REsourceExists"))

	http.Handle("PUT /", otelhttp.NewHandler(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		handler.ReplaceResource(resp, req, repository)
	}), "ReplcaeResource"))

	http.Handle("DELETE /", otelhttp.NewHandler(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		handler.ReplaceResource(resp, req, repository)
	}), "RemoveResource"))
}

func Run(host string, port int) {
	http.ListenAndServe(host+":"+strconv.Itoa(port), nil)
}
