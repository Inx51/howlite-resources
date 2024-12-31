package handler

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/inx51/howlite/resources/api/resource"
	reserr "github.com/inx51/howlite/resources/api/resource/errors"
)

func GetResource(resp *http.ResponseWriter, req *http.Request) {
	identifier := resource.GetIdentifier(&req.URL.Path)
	res, err := resource.Get(&identifier)
	if err != nil {
		if errors.Is(err, reserr.NotFoundError{Identifier: identifier}) {
			slog.Warn("Failed to get resource.", slog.Any("error", err), slog.Any("identifier", identifier))
			(*resp).WriteHeader(404)
			return
		}
		slog.Error("Unhandled error.", slog.Any("error", err), slog.Any("identifier", identifier))
		(*resp).WriteHeader(500)
		return
	}
	(*resp).WriteHeader(201)
	defer res.Body.Close()
	buff := make([]byte, 1024)
	io.CopyBuffer((*resp), res.Body, buff)
}
