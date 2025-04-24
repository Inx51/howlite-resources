package handler

import (
	"log/slog"
	"net/http"

	"github.com/inx51/howlite/resources/resource"
	"github.com/inx51/howlite/resources/resource/repository"
	"go.opentelemetry.io/otel/sdk/metric"
)

func RemoveResource(
	resp http.ResponseWriter,
	req *http.Request,
	repository *repository.Repository,
	logger *slog.Logger,
	meter *metric.MeterProvider) error {
	resourceIdentifier := resource.NewResourceIdentifier(&req.URL.Path)

	resourceExists, err := repository.ResourceExists(resourceIdentifier)
	if err != nil {
		resp.WriteHeader(500)
		return err
	}

	if !resourceExists {
		logger.Debug("Failed to remove resource since it does not exist", "resourceIdentifier", resourceIdentifier.Value)
		resp.WriteHeader(404)
		return nil
	}

	err = repository.RemoveResource(resourceIdentifier)
	if err != nil {
		return err
	}

	resp.WriteHeader(204)
	logger.Info("Removed resource", "resourceIdentifier", resourceIdentifier.Value)
	return nil
}
