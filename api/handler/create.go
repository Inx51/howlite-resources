package handler

import (
	"net/http"

	"github.com/inx51/howlite/resources/resource"
	"github.com/inx51/howlite/resources/resource/repository"
)

func CreateResource(resp http.ResponseWriter, req *http.Request, repository *repository.Repository) {
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

	resp.Header().Add("Location", req.URL.Scheme+"://"+req.URL.Host+req.URL.Path)
	resp.WriteHeader(201)
}
