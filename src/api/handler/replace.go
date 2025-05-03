package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/inx51/howlite/resources/api/handler/services"
	"github.com/inx51/howlite/resources/resource"
	"github.com/inx51/howlite/resources/resource/repository"
	"go.opentelemetry.io/otel/sdk/metric"
)

func ReplaceResource(
	ctx context.Context,
	resp http.ResponseWriter,
	req *http.Request,
	repository *repository.Repository,
	logger *slog.Logger,
	meter *metric.MeterProvider) error {
	resourceIdentifier := resource.NewResourceIdentifier(&req.URL.Path)

	resourceExists, err := repository.ResourceExistsContext(ctx, resourceIdentifier)
	if err != nil {
		return err
	}

	headers := make(map[string][]string)
	for k, v := range req.Header {
		headers[k] = v
	}

	resource := resource.NewResource(resourceIdentifier, &headers, &req.Body)
	err = repository.SaveResourceContext(ctx, resource)
	if err != nil {
		return err
	}

	location := services.GetRequestUrl(req)
	resp.Header().Add("Location", location)
	if !resourceExists {
		logger.InfoContext(ctx, "Resource created", "resourceIdentifier", resourceIdentifier.Value)
		resp.WriteHeader(201)
	} else {
		logger.InfoContext(ctx, "Existing resource replaced", "resourceIdentifier", resourceIdentifier.Value)
		resp.WriteHeader(204)
	}
	return nil
}
