package handlers

import (
	"context"
	"net/http"

	"github.com/inx51/howlite-resources/http/uri"
	"github.com/inx51/howlite-resources/logger"
	"github.com/inx51/howlite-resources/meter"
	"github.com/inx51/howlite-resources/resource"
	"github.com/inx51/howlite-resources/storage"
	"github.com/inx51/howlite-resources/tracer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type ReplaceHandler struct {
	storage *storage.Storage
}

func (handler *ReplaceHandler) Method() string {
	return "PUT"
}

func (handler *ReplaceHandler) Path() string {
	return "/"
}

func (handler *ReplaceHandler) Handle(
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

	headers := make(map[string][]string)
	for k, v := range req.Header {
		headers[k] = v
	}

	resource := resource.NewResource(resourceIdentifier, &req.Body)
	defer (*resource.Body).Close()
	for k, v := range req.Header {
		resource.Headers.Add(ctx, k, v)
	}
	srCtx, span := tracer.StartInfoSpan(ctx, "storage."+storage.GetName()+".save_resource")
	tracer.SetInfoAttributes(
		srCtx,
		span,
		attribute.String("resource_identifier", resourceIdentifier.Identifier()),
	)
	err = storage.SaveResource(srCtx, resource)
	tracer.SafeEndSpan(span)
	if err != nil {
		statusCode = http.StatusInternalServerError
		resp.WriteHeader(statusCode)
		return statusCode, err
	}

	location := uri.AbsoluteUri(req)
	resp.Header().Add("Location", location)
	if !resourceExists {
		meter.ArithmeticInt64Counter(ctx, "resources_created_total", 1, metric.WithAttributes(attribute.String("resource_identifier", resourceIdentifier.Identifier())))
		meter.ArithmeticInt64Counter(ctx, "resources_overall", 1)
		logger.Info(ctx, "Resource created", "resourceIdentifier", resourceIdentifier.Identifier())
		statusCode = http.StatusCreated
		resp.WriteHeader(statusCode)
	} else {
		meter.ArithmeticInt64Counter(ctx, "resources_replaced_total", 1, metric.WithAttributes(attribute.String("resource_identifier", resourceIdentifier.Identifier())))
		logger.Info(ctx, "Existing resource replaced", "resourceIdentifier", resourceIdentifier.Identifier())
		resp.WriteHeader(statusCode)
	}
	return statusCode, nil
}

func NewReplaceHandler(storage *storage.Storage) Handler {
	return &ReplaceHandler{
		storage: storage,
	}
}
