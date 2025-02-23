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
	logger *slog.Logger) {
	resourceIdentifier := resource.NewResourceIdentifier(&req.URL.Path)

	resourceExists, existsErr := repository.ResourceExists(resourceIdentifier)
	if existsErr != nil {
		resp.WriteHeader(500)
		panic(existsErr)
	}

	if resourceExists {
		resp.WriteHeader(409)
		return
	}

	headers := make(map[string][]string)
	for k, v := range req.Header {
		headers[k] = v
	}

	resource := resource.NewResource(resourceIdentifier, &headers, &req.Body)
	err := repository.SaveResource(resource)
	if err != nil {
		panic(err)
	}

	location := services.GetRequestUrl(req)
	resp.Header().Add("Location", location)
	resp.WriteHeader(201)
}
