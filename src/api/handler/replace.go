package handler

import (
	"log/slog"
	"net/http"

	"github.com/inx51/howlite/resources/api/handler/services"
	"github.com/inx51/howlite/resources/resource"
	"github.com/inx51/howlite/resources/resource/repository"
)

func ReplaceResource(
	resp http.ResponseWriter,
	req *http.Request,
	repository *repository.Repository,
	logger *slog.Logger) {
	resourceIdentifier := resource.NewResourceIdentifier(&req.URL.Path)

	resourceExists, existsErr := repository.ResourceExists(resourceIdentifier)
	if existsErr != nil {
		panic(existsErr)
	}

	headers := make(map[string][]string)
	for k, v := range req.Header {
		headers[k] = v
	}

	resource := resource.NewResource(resourceIdentifier, &headers, &req.Body)
	saveErr := repository.SaveResource(resource)
	if saveErr != nil {
		panic(saveErr)
	}

	location := services.GetRequestUrl(req)
	resp.Header().Add("Location", location)
	if !resourceExists {
		resp.WriteHeader(201)
	} else {
		resp.WriteHeader(204)
	}
}
