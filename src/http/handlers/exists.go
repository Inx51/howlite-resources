package handlers

import (
	"context"
	"net/http"

	"github.com/inx51/howlite-resources/http/response"
	"github.com/inx51/howlite-resources/logger"
	"github.com/inx51/howlite-resources/resource"
	"github.com/inx51/howlite-resources/storage"
	"github.com/inx51/howlite-resources/tracer"
	"go.opentelemetry.io/otel/attribute"
)

type ExistsHandler struct {
	storage *storage.Storage
}

func (handler *ExistsHandler) Method() string {
	return "HEAD"
}

func (handler *ExistsHandler) Path() string {
	return "/"
}

func (handler *ExistsHandler) Handle(
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
	exists, err := storage.ResourceExists(reCtx, resourceIdentifier)
	tracer.SafeEndSpan(span)
	if err != nil {
		statusCode = http.StatusInternalServerError
		resp.WriteHeader(statusCode)
		return statusCode, err
	}

	if exists {

		grCtx, span := tracer.StartInfoSpan(ctx, "storage."+storage.GetName()+".get_resource")
		tracer.SetInfoAttributes(
			grCtx,
			span,
			attribute.String("resource_identifier", resourceIdentifier.Identifier()),
		)
		resource, _ := storage.GetResource(grCtx, resourceIdentifier)
		span.End()
		defer (*resource.Body).Close()

		response.WriteHeaders(resource.Headers.Headers(), resp)

		logger.Debug(ctx, "Resource found", "resourceIdentifier", resourceIdentifier.Identifier())
		resp.WriteHeader(statusCode)
	} else {
		logger.Debug(ctx, "Failed to find resource", "resourceIdentifier", resourceIdentifier.Identifier())
		statusCode = http.StatusNoContent
		resp.WriteHeader(statusCode)
	}
	return statusCode, nil
}

func NewExistsHandler(storage *storage.Storage) Handler {
	return &ExistsHandler{
		storage: storage,
	}
}
