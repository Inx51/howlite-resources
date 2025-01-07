package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/inx51/howlite/resources/resource"
	"github.com/inx51/howlite/resources/resource/service"
	"github.com/inx51/howlite/resources/storage"
)

func RemoveResource(
	resp http.ResponseWriter,
	req *http.Request,
	storage *storage.Storage,
) {
	identifier := resource.NewIdentifier(&req.URL.Path)
	err := service.Remove(&identifier, storage)
	if errors.Is(err, resource.NotFoundError{Identifier: &identifier}) {
		slog.Warn("Failed to delete resource since it doesnt exist.", slog.Any("error", err), slog.Any("identifier", identifier.Value))
		resp.WriteHeader(404)
		return
	}
	resp.WriteHeader(204)
}
