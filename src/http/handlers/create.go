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

type CreateHandler struct {
	storage *storage.Storage
}

func (handler *CreateHandler) Method() string {
	return "POST"
}

func (handler *CreateHandler) Path() string {
	return "/"
}

func (handler *CreateHandler) Handle(
	ctx context.Context,
	req *http.Request,
	resp http.ResponseWriter) (int, error) {
	resourceIdentifier := resource.NewResourceIdentifier(req.URL.Path)

	storage := *handler.storage
	statusCode := http.StatusCreated
	reCtx, span := tracer.StartInfoSpan(ctx, "storage."+storage.GetName()+".resource_exists")
	tracer.SetInfoAttributes(
		reCtx,
		span,
		attribute.String("resource.identifier", resourceIdentifier.Identifier()),
	)
	resourceExists, err := storage.ResourceExists(reCtx, resourceIdentifier)
	tracer.SafeEndSpan(span)
	if err != nil {
		statusCode = http.StatusInternalServerError
		resp.WriteHeader(statusCode)
		return statusCode, err
	}

	if resourceExists {
		logger.Debug(ctx, "Can't create resource with the given identifier because it already exists", "resourceIdentifier", resourceIdentifier.Identifier())
		statusCode = http.StatusConflict
		resp.WriteHeader(statusCode)
		return statusCode, nil
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

	meter.ArithmeticInt64Counter(ctx, "resources_created_total", 1, metric.WithAttributes(attribute.String("resource_identifier", resourceIdentifier.Identifier())))
	meter.ArithmeticInt64Counter(ctx, "resources_overall", 1)

	location := uri.AbsoluteUri(req)
	resp.Header().Add("Location", location)
	resp.WriteHeader(statusCode)
	logger.Info(ctx, "Resource created", "resourceIdentifier", resourceIdentifier.Identifier())
	return statusCode, nil
}

func NewCreateHandler(storage *storage.Storage) Handler {
	return &CreateHandler{
		storage: storage,
	}
}
