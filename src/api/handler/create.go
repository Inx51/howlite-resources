package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/inx51/howlite/resources/resource"
	"github.com/inx51/howlite/resources/url"
)

func CreateResource(resp *http.ResponseWriter, req *http.Request) {
	identifier := resource.NewIdentifier(&req.URL.Path)
	res := resource.New(&identifier, req.Header, &req.Body)
	err := resource.Create(&res)
	if err != nil {
		if errors.Is(err, resource.AlreadyExistsError{Identifier: &identifier}) {
			slog.Warn("Failed to create resource.", slog.Any("error", err), slog.Any("identifier", identifier.Value))
			(*resp).WriteHeader(409)
			return
		}
		slog.Error("Unhandled error.", slog.Any("error", err), slog.Any("identifier", identifier.Value))
		return
	}
	(*resp).WriteHeader(201)
	(*resp).Header().Add("Location", url.GetAbsolute(req.URL))
}
