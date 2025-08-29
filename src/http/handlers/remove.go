package handlers

import (
	"context"
	"net/http"

	"github.com/inx51/howlite-resources/logger"
	"github.com/inx51/howlite-resources/meter"
	"github.com/inx51/howlite-resources/resource"
	"github.com/inx51/howlite-resources/storage"
	"github.com/inx51/howlite-resources/tracer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type RemoveHandler struct {
	storage *storage.Storage
}

func (handler *RemoveHandler) Method() string {
	return "DELETE"
}

func (handler *RemoveHandler) Path() string {
	return "/"
}

func (handler *RemoveHandler) Handle(
	ctx context.Context,
	req *http.Request,
	resp http.ResponseWriter) (int, error) {
	resourceIdentifier := resource.NewResourceIdentifier(req.URL.Path)

	storage := *handler.storage
	statusCode := http.StatusNoContent
	reCtx, span := tracer.StartInfoSpan(ctx, "storage."+storage.GetName()+".resource_exists")
	tracer.SetInfoAttributes(
		reCtx,
		span,
		attribute.String("resource_identifier", resourceIdentifier.Identifier()),
	)
	resourceExists, err := storage.ResourceExists(reCtx, resourceIdentifier)
	tracer.SafeEndSpan(span)
	if err != nil {
		statusCode = http.StatusInternalServerError
		resp.WriteHeader(statusCode)
		return statusCode, err
	}

	if !resourceExists {
		logger.Debug(ctx, "Failed to get resource since it does not exist", "resourceIdentifier", resourceIdentifier.Identifier())
		statusCode = http.StatusNotFound
		resp.WriteHeader(statusCode)
		return statusCode, nil
	}

	rrCtx, span := tracer.StartInfoSpan(ctx, "storage."+storage.GetName()+".remove_resource")
	tracer.SetInfoAttributes(
		rrCtx,
		span,
		attribute.String("resource_identifier", resourceIdentifier.Identifier()),
	)
	err = storage.RemoveResource(rrCtx, resourceIdentifier)
	tracer.SafeEndSpan(span)
	if err != nil {
		statusCode = http.StatusInternalServerError
		return statusCode, err
	}

	meter.ArithmeticInt64Counter(ctx, "resources_removed_total", 1, metric.WithAttributes(attribute.String("resource_identifier", resourceIdentifier.Identifier())))

	resp.WriteHeader(statusCode)
	logger.Info(ctx, "Removed resource", "resourceIdentifier", resourceIdentifier.Identifier())
	return statusCode, nil
}

func NewRemoveHandler(storage *storage.Storage) Handler {
	return &RemoveHandler{
		storage: storage,
	}
}
