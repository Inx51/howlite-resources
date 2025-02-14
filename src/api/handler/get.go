package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/inx51/howlite/resources/resource"
	"github.com/inx51/howlite/resources/resource/repository"
)

func GetResource(resp http.ResponseWriter, req *http.Request, repository *repository.Repository) {
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

	resource, err := repository.GetResource(resourceIdentifier)
	if err != nil {
		resp.WriteHeader(500)
		panic(err)
	}

	for k, v := range *resource.Headers {
		resp.Header().Add(k, strings.Join(v, ",'"))
	}

	resp.WriteHeader(200)

	buff := make([]byte, 1024)
	body := *resource.Body
	io.CopyBuffer(resp, body, buff)
	body.Close()
}
