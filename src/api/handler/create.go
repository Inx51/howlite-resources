package handler

import (
	"log/slog"
	"net/http"

	"github.com/inx51/howlite/resources/api/handler/services"
	"github.com/inx51/howlite/resources/resource"
	"github.com/inx51/howlite/resources/resource/repository"
)

func CreateResource(
	resp http.ResponseWriter,
	req *http.Request,
	repository *repository.Repository,
	logger *slog.Logger) error {
	resourceIdentifier := resource.NewResourceIdentifier(&req.URL.Path)

	resourceExists, err := repository.ResourceExists(resourceIdentifier)
	if err != nil {
		resp.WriteHeader(500)
		return err
	}

	if resourceExists {
		logger.Debug("Can't create resource with the given identifier because it already exists", "resourceIdentifier", resourceIdentifier.Value)
		resp.WriteHeader(409)
		return nil
	}

	headers := make(map[string][]string)
	for k, v := range req.Header {
		headers[k] = v
	}

	resource := resource.NewResource(resourceIdentifier, &headers, &req.Body)
	err = repository.SaveResource(resource)
	if err != nil {
		resp.WriteHeader(500)
		return err
	}

	location := services.GetRequestUrl(req)
	resp.Header().Add("Location", location)
	resp.WriteHeader(201)
	logger.Debug("Resource created", "resourceIdentifier", resourceIdentifier.Value)
	return nil
}
