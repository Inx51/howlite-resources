package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/inx51/howlite/resources/api/resource"
	reserr "github.com/inx51/howlite/resources/api/resource/errors"
	"github.com/inx51/howlite/resources/api/url"
)

func CreateResource(resp *http.ResponseWriter, req *http.Request) {
	identifier := resource.GetIdentifier(&req.URL.Path)
	err := resource.Create(&identifier, &req.Body, &req.Header)
	if err != nil {
		if errors.Is(err, reserr.AlreadyExistsError{Identifier: identifier}) {
			slog.Warn("Failed to create resource.", slog.Any("error", err), slog.Any("identifier", identifier))
			(*resp).WriteHeader(409)
			return
		}
		slog.Error("Unhandled error.", slog.Any("error", err), slog.Any("identifier", identifier))
		return
	}
	(*resp).WriteHeader(201)
	(*resp).Header().Add("Location", url.GetAbsolute(req.URL))
}
