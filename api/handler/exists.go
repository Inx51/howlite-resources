package handler

import (
	"net/http"

	"github.com/inx51/howlite/resources/resource"
	"github.com/inx51/howlite/resources/resource/repository"
)

func ResourceExists(resp http.ResponseWriter, req *http.Request, repository *repository.Repository) {
	resourceIdentifier := resource.NewResourceIdentifier(&req.URL.Path)
	exists, err := repository.ResourceExists(resourceIdentifier)
	if err != nil {
		resp.WriteHeader(500)
		panic(err)
	}

	if exists {
		resp.WriteHeader(200)
	} else {
		resp.WriteHeader(404)
	}
}
