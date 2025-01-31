package handler

import (
	"net/http"

	"github.com/inx51/howlite/resources/resource"
	"github.com/inx51/howlite/resources/resource/repository"
)

func ReplaceResource(resp http.ResponseWriter, req *http.Request, repository *repository.Repository) {
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
	defer (*resource.Body).Close()
	saveErr := repository.SaveResource(resource)
	if saveErr != nil {
		panic(saveErr)
	}

	if !resourceExists {
		resp.WriteHeader(201)
		resp.Header().Add("Location", req.URL.Scheme+"://"+req.URL.Host+req.URL.Path)
		return
	} else {
		resp.WriteHeader(204)
		resp.Header().Add("Location", req.URL.Scheme+"://"+req.URL.Host+req.URL.Path)
	}
}
