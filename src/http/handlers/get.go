package handlers

import (
	"context"
	"net/http"

	"github.com/inx51/howlite-resources/http/response"
	"github.com/inx51/howlite-resources/logger"
	"github.com/inx51/howlite-resources/meter"
	"github.com/inx51/howlite-resources/resource"
	"github.com/inx51/howlite-resources/storage"
	"github.com/inx51/howlite-resources/tracer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type GetHandler struct {
	storage *storage.Storage
	buffer  []byte
}

func (handler *GetHandler) Method() string {
	return "GET"
}

func (handler *GetHandler) Path() string {
	return "/"
}

func (handler *GetHandler) Handle(
	ctx context.Context,
	req *http.Request,
	resp http.ResponseWriter) (int, error) {
	resourceIdentifier := resource.NewResourceIdentifier(req.URL.Path)

	storage := *handler.storage
	statusCode := 200
	reCtx, span := tracer.StartInfoSpan(ctx, "storage."+storage.GetName()+".resource_exists")
	tracer.SetInfoAttributes(
		reCtx,
		span,
		attribute.String("resource_identifier", resourceIdentifier.Identifier()),
	)
	resourceExists, err := storage.ResourceExists(reCtx, resourceIdentifier)
	tracer.SafeEndSpan(span)
	if err != nil {
		statusCode = 500
		resp.WriteHeader(statusCode)
		return statusCode, err
	}

	if !resourceExists {
		logger.Debug(ctx, "Failed to get resource since it does not exist", "resourceIdentifier", resourceIdentifier.Identifier())
		statusCode = 404
		resp.WriteHeader(statusCode)
		return statusCode, nil
	}

	grCtx, span := tracer.StartInfoSpan(ctx, "storage."+storage.GetName()+".get_resource")
	resource, err := storage.GetResource(grCtx, resourceIdentifier)
	tracer.SetInfoAttributes(
		grCtx,
		span,
		attribute.String("resource_identifier", resourceIdentifier.Identifier()),
	)
	tracer.SafeEndSpan(span)

	response.WriteHeaders(resource.Headers.Headers(), resp)

	if err != nil {
		statusCode = 500
		resp.WriteHeader(statusCode)
		return statusCode, err
	}

	meter.ArithmeticInt64Counter(ctx, "resources_fetched_total", 1, metric.WithAttributes(attribute.String("resource_identifier", resourceIdentifier.Identifier())))

	resp.WriteHeader(statusCode)

	response.WriteBody(*resource.Body, resp)
	logger.Debug(ctx, "Resource returned", "resourceIdentifier", resourceIdentifier.Identifier())
	return statusCode, nil
}

func NewGetHandler(storage *storage.Storage) Handler {
	return &GetHandler{
		storage: storage,
		buffer:  make([]byte, 1024),
	}
}
