package handler

import (
	"net/http"
	"strings"

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

		resource, _ := repository.GetResource(resourceIdentifier)
		defer (*resource.Body).Close()
		for k, v := range *resource.Headers {
			resp.Header().Add(k, strings.Join(v, ",'"))
		}

		resp.WriteHeader(204)
	} else {
		resp.WriteHeader(404)
	}
}
