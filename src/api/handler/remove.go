package handler

import (
	"context"
	"log/slog"
	"net/http"
	"sync"

	"github.com/inx51/howlite/resources/resource"
	"github.com/inx51/howlite/resources/resource/repository"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	resourcesRemovedCounter     metric.Int64Counter
	resourcesRemovedCounterOnce sync.Once
)

func RemoveResource(
	ctx context.Context,
	resp http.ResponseWriter,
	req *http.Request,
	repository *repository.Repository,
	logger *slog.Logger,
	meter *metric.Meter) error {
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

	incrementResourcesRemovedCounterContext(ctx, meter, resourceIdentifier)

	resp.WriteHeader(204)
	logger.InfoContext(ctx, "Removed resource", "resourceIdentifier", resourceIdentifier.Value)
	return nil
}

func incrementResourcesRemovedCounterContext(ctx context.Context, meter *metric.Meter, resourceIdentifier *resource.ResourceIdentifier) {
	resourcesRemovedCounterOnce.Do(func() {
		resourcesRemovedCounter, _ = (*meter).Int64Counter("resources.removed_total")
	})
	resourcesRemovedCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("resource_identifier", *resourceIdentifier.Value)))
}
