package api

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/inx51/howlite/resources/api/handler"
	"github.com/inx51/howlite/resources/resource/repository"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/sdk/metric"
)

type Endpoint struct {
	Method             string
	Description        string
	HandlerWithContext Handler
}

type Handler func(context.Context, http.ResponseWriter, *http.Request, *repository.Repository, *slog.Logger, *metric.MeterProvider) error

func SetupHandlers(
	repository *repository.Repository,
	logger *slog.Logger,
	meter *metric.MeterProvider) {

	endpoints := []Endpoint{
		{
			Method:             "GET",
			Description:        "GetResource",
			HandlerWithContext: handler.GetResource,
		},
		{
			Method:             "POST",
			Description:        "CreateResource",
			HandlerWithContext: handler.CreateResource,
		},
		{
			Method:             "HEAD",
			Description:        "ResourceExists",
			HandlerWithContext: handler.ResourceExists,
		},
		{
			Method:             "PUT",
			Description:        "ReplaceResource",
			HandlerWithContext: handler.ReplaceResource,
		},
		{
			Method:             "DELETE",
			Description:        "RemoveResource",
			HandlerWithContext: handler.RemoveResource,
		},
	}

	http.DefaultServeMux = http.NewServeMux()

	for _, endpoint := range endpoints {
		http.Handle(endpoint.Method+" /", otelhttp.NewHandler(
			http.HandlerFunc(
				func(resp http.ResponseWriter, req *http.Request) {
					// Extract the span from the request's context
					ctx := req.Context()

					logger.InfoContext(ctx, "Request received", "method", req.Method, "url", req.URL.Path)
					logger.DebugContext(ctx, "Found matching endpoint route", "method", endpoint.Method, "path", req.URL.Path)
					err := endpoint.HandlerWithContext(
						ctx,
						resp,
						req,
						repository,
						logger,
						meter)
					if err != nil {
						logger.ErrorContext(ctx, "Unhandled error occurred", "error", err)
					} else {
						logger.InfoContext(ctx, "Response sent", "method", req.Method, "url", req.URL.Path)
					}
				},
			),
			endpoint.Description),
		)
	}
}

func RunWithContext(
	ctx context.Context,
	host string,
	port int,
	logger *slog.Logger) {
	logger.InfoContext(ctx, "Starting HTTP server", "host", host, "port", port)
	http.ListenAndServe(host+":"+strconv.Itoa(port), nil)
}
