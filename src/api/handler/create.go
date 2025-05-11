package handler

import (
	"context"
	"log/slog"
	"net/http"
	"sync"

	"github.com/inx51/howlite/resources/api/handler/services"
	"github.com/inx51/howlite/resources/resource"
	"github.com/inx51/howlite/resources/resource/repository"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	resourcesCreatedCounter     metric.Int64Counter
	resourcesCreatedCounterOnce sync.Once
)

func CreateResource(
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

	if resourceExists {
		logger.DebugContext(ctx, "Can't create resource with the given identifier because it already exists", "resourceIdentifier", resourceIdentifier.Value)
		resp.WriteHeader(409)
		return nil
	}

	headers := make(map[string][]string)
	for k, v := range req.Header {
		headers[k] = v
	}

	resource := resource.NewResource(resourceIdentifier, &headers, &req.Body)
	err = repository.SaveResourceContext(ctx, resource)
	if err != nil {
		resp.WriteHeader(500)
		return err
	}

	incrementResourceCreatedCounterContext(ctx, meter, resourceIdentifier)

	location := services.GetRequestUrl(req)
	resp.Header().Add("Location", location)
	resp.WriteHeader(201)
	logger.InfoContext(ctx, "Resource created", "resourceIdentifier", resourceIdentifier.Value)
	return nil
}

func incrementResourceCreatedCounterContext(ctx context.Context, meter *metric.Meter, resourceIdentifier *resource.ResourceIdentifier) {
	resourcesCreatedCounterOnce.Do(func() {
		resourcesCreatedCounter, _ = (*meter).Int64Counter("resources.created_total")
	})
	resourcesCreatedCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("resource_identifier", *resourceIdentifier.Value)))
}
