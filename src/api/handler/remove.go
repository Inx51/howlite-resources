package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/inx51/howlite/resources/resource"
	"github.com/inx51/howlite/resources/resource/repository"
	"go.opentelemetry.io/otel/sdk/metric"
)

func RemoveResource(
	ctx context.Context,
	resp http.ResponseWriter,
	req *http.Request,
	repository *repository.Repository,
	logger *slog.Logger,
	meter *metric.MeterProvider) error {
	resourceIdentifier := resource.NewResourceIdentifier(&req.URL.Path)

	resourceExists, err := repository.ResourceExistsContext(ctx, resourceIdentifier)
	if err != nil {
		resp.WriteHeader(500)
		return err
	}

	if !resourceExists {
		logger.DebugContext(ctx, "Failed to remove resource since it does not exist", "resourceIdentifier", resourceIdentifier.Value)
		resp.WriteHeader(404)
		return nil
	}

	err = repository.RemoveResourceContext(ctx, resourceIdentifier)
	if err != nil {
		return err
	}

	resp.WriteHeader(204)
	logger.InfoContext(ctx, "Removed resource", "resourceIdentifier", resourceIdentifier.Value)
	return nil
}
