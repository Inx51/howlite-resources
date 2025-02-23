package handler

import (
	"log/slog"
	"net/http"

	"github.com/inx51/howlite/resources/resource"
	"github.com/inx51/howlite/resources/resource/repository"
)

func RemoveResource(
	resp http.ResponseWriter,
	req *http.Request,
	repository *repository.Repository,
	logger *slog.Logger) {
	resourceIdentifier := resource.NewResourceIdentifier(&req.URL.Path)

	resourceExists, existsErr := repository.ResourceExists(resourceIdentifier)
	if existsErr != nil {
		resp.WriteHeader(500)
		panic(existsErr)
	}

	if !resourceExists {
		resp.WriteHeader(404)
		return
	}

	err := repository.RemoveResource(resourceIdentifier)
	if err != nil {
		panic(err)
	}

	resp.WriteHeader(204)
}
