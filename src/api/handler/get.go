package handler

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/inx51/howlite/resources/resource"
	"github.com/inx51/howlite/resources/resource/repository"
	"go.opentelemetry.io/otel/sdk/metric"
)

func GetResource(
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
		logger.DebugContext(ctx, "Failed to get resource since it does not exist", "resourceIdentifier", resourceIdentifier.Value)
		resp.WriteHeader(404)
		return nil
	}

	resource, err := repository.GetResourceContext(ctx, resourceIdentifier)
	if err != nil {
		resp.WriteHeader(500)
		return err
	}

	for k, v := range *resource.Headers {
		resp.Header().Add(k, strings.Join(v, ",'"))
	}

	resp.WriteHeader(200)

	buff := make([]byte, 1024)
	body := *resource.Body
	io.CopyBuffer(resp, body, buff)
	body.Close()
	logger.DebugContext(ctx, "Resource returned", "resourceIdentifier", resourceIdentifier.Value)
	return nil
}
