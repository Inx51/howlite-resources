package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/inx51/howlite/resources/resource"
)

func RemoveResource(resp *http.ResponseWriter, req *http.Request) {
	identifier := resource.NewIdentifier(&req.URL.Path)
	err := resource.Remove(&identifier)
	if errors.Is(err, resource.NotFoundError{Identifier: &identifier}) {
		slog.Warn("Failed to delete resource since it doesnt exist.", slog.Any("error", err), slog.Any("identifier", identifier.Value))
		(*resp).WriteHeader(404)
		return
	}
	(*resp).WriteHeader(204)
}
