package api

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/inx51/howlite/resources/api/handler"
	"github.com/inx51/howlite/resources/resource/repository"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/sdk/metric"
)

type Endpoint struct {
	Method      string
	Description string
	Handler     Handler
}

type Handler func(http.ResponseWriter, *http.Request, *repository.Repository, *slog.Logger, *metric.MeterProvider) error

func SetupHandlers(
	repository *repository.Repository,
	logger *slog.Logger,
	meter *metric.MeterProvider) {

	endpoints := []Endpoint{
		{
			Method:      "GET",
			Description: "GetResource",
			Handler:     handler.GetResource,
		},
		{
			Method:      "POST",
			Description: "CreateResource",
			Handler:     handler.CreateResource,
		},
		{
			Method:      "HEAD",
			Description: "ResourceExists",
			Handler:     handler.ResourceExists,
		},
		{
			Method:      "PUT",
			Description: "ReplaceResource",
			Handler:     handler.ReplaceResource,
		},
		{
			Method:      "DELETE",
			Description: "RemoveResource",
			Handler:     handler.RemoveResource,
		},
	}

	http.DefaultServeMux = http.NewServeMux()

	for _, endpoint := range endpoints {
		http.Handle(endpoint.Method+" /", otelhttp.NewHandler(
			http.HandlerFunc(
				func(resp http.ResponseWriter, req *http.Request) {
					logger.Info("Request received", "method", req.Method, "url", req.URL.Path)
					logger.Debug("Found matching endpoint route", "method", endpoint.Method, "path", req.URL.Path)
					err := endpoint.Handler(
						resp,
						req,
						repository,
						logger,
						meter)
					if err != nil {
						logger.Error("Unhandled error occurred", "error", err)
					} else {
						logger.Info("Response sent", "method", req.Method, "url", req.URL.Path)
					}
				},
			),
			endpoint.Description),
		)
	}
}

func Run(
	host string,
	port int,
	logger *slog.Logger) {
	logger.Info("Starting HTTP server", "host", host, "port", port)
	http.ListenAndServe(host+":"+strconv.Itoa(port), nil)
}
