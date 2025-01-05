package handler

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/inx51/howlite/resources/resource"
)

func GetResource(resp *http.ResponseWriter, req *http.Request) {
	identifier := resource.NewIdentifier(&req.URL.Path)
	res, err := resource.Get(&identifier)
	if err != nil {
		if errors.Is(err, resource.NotFoundError{Identifier: &identifier}) {
			slog.Warn("Failed to get resource.", slog.Any("error", err), slog.Any("identifier", identifier.Value))
			(*resp).WriteHeader(404)
			return
		}
		slog.Error("Unhandled error.", slog.Any("error", err), slog.Any("identifier", identifier.Value))
		(*resp).WriteHeader(500)
		return
	}
	if len(res.Headers) > 0 {
		for k, v := range res.Headers {
			(*resp).Header().Add(k, strings.Join(v, ",'"))
		}
	}
	(*resp).WriteHeader(201)
	defer (*res.Body).Close()

	buff := make([]byte, 1024)
	io.CopyBuffer((*resp), *res.Body, buff)
}
