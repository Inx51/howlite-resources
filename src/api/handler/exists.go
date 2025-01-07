package handler

import (
	"net/http"

	"github.com/inx51/howlite/resources/resource"
	"github.com/inx51/howlite/resources/resource/service"
	"github.com/inx51/howlite/resources/storage"
)

func ResourceExists(
	resp http.ResponseWriter,
	req *http.Request,
	storage *storage.Storage,
) {
	identifier := resource.NewIdentifier(&req.URL.Path)
	exist := service.Exists(&identifier, storage)
	if exist {
		resp.WriteHeader(204)
	} else {
		resp.WriteHeader(404)
	}
}
