package handler

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/inx51/howlite/resources/resource"
	"github.com/inx51/howlite/resources/resource/repository"
)

func ResourceExists(
	resp http.ResponseWriter,
	req *http.Request,
	repository *repository.Repository,
	logger *slog.Logger) error {
	resourceIdentifier := resource.NewResourceIdentifier(&req.URL.Path)
	exists, err := repository.ResourceExists(resourceIdentifier)
	if err != nil {
		resp.WriteHeader(500)
		return err
	}

	if exists {

		resource, _ := repository.GetResource(resourceIdentifier)
		defer (*resource.Body).Close()
		for k, v := range *resource.Headers {
			resp.Header().Add(k, strings.Join(v, ",'"))
		}
		logger.Debug("Resource found", "resourceIdentifier", resourceIdentifier.Value)
		resp.WriteHeader(204)
	} else {
		logger.Debug("Failed to find resource", "resourceIdentifier", resourceIdentifier.Value)
		resp.WriteHeader(404)
	}
	return nil
}
