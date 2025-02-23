package api

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/inx51/howlite/resources/api/handler"
	"github.com/inx51/howlite/resources/resource/repository"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func SetupHandlers(repository *repository.Repository, logger *slog.Logger) {
	http.DefaultServeMux = http.NewServeMux()

	http.Handle("GET /", otelhttp.NewHandler(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		handler.GetResource(resp, req, repository, logger)
	}), "GetResource"))

	http.Handle("POST /", otelhttp.NewHandler(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		handler.CreateResource(resp, req, repository, logger)
	}), "CreateResource"))

	http.Handle("HEAD /", otelhttp.NewHandler(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		handler.ResourceExists(resp, req, repository, logger)
	}), "REsourceExists"))

	http.Handle("PUT /", otelhttp.NewHandler(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		handler.ReplaceResource(resp, req, repository, logger)
	}), "ReplcaeResource"))

	http.Handle("DELETE /", otelhttp.NewHandler(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		handler.RemoveResource(resp, req, repository, logger)
	}), "RemoveResource"))
}

func Run(
	host string,
	port int,
	logger *slog.Logger) {
	logger.Info("Starting HTTP server")
	http.ListenAndServe(host+":"+strconv.Itoa(port), nil)
	logger.Info("Now listening", "host", host, "port", port)
}
