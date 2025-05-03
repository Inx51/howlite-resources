package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/inx51/howlite/resources/resource"
	"github.com/inx51/howlite/resources/resource/repository"
	"go.opentelemetry.io/otel/sdk/metric"
)

func ResourceExists(
	ctx context.Context,
	resp http.ResponseWriter,
	req *http.Request,
	repository *repository.Repository,
	logger *slog.Logger,
	meter *metric.MeterProvider) error {
	resourceIdentifier := resource.NewResourceIdentifier(&req.URL.Path)
	exists, err := repository.ResourceExistsContext(ctx, resourceIdentifier)
	if err != nil {
		resp.WriteHeader(500)
		return err
	}

	if exists {

		resource, _ := repository.GetResourceContext(ctx, resourceIdentifier)
		defer (*resource.Body).Close()
		for k, v := range *resource.Headers {
			resp.Header().Add(k, strings.Join(v, ",'"))
		}
		logger.DebugContext(ctx, "Resource found", "resourceIdentifier", resourceIdentifier.Value)
		resp.WriteHeader(204)
	} else {
		logger.DebugContext(ctx, "Failed to find resource", "resourceIdentifier", resourceIdentifier.Value)
		resp.WriteHeader(404)
	}
	return nil
}
