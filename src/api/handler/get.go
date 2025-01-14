package handler

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/inx51/howlite/resources/resource"
	"github.com/inx51/howlite/resources/resource/service"
	"github.com/inx51/howlite/resources/storage"
)

func GetResource(
	resp http.ResponseWriter,
	req *http.Request,
	storage *storage.Storage,
) {
	identifier := resource.NewIdentifier(&req.URL.Path)
	res, err := service.Get(&identifier, storage)
	if err != nil {
		if errors.Is(err, resource.NotFoundError{Identifier: &identifier}) {
			slog.Warn("Failed to get resource since it doesnt exist.", slog.Any("error", err), slog.Any("identifier", identifier.Value))
			resp.WriteHeader(404)
			return
		}
		slog.Error("Unhandled error.", slog.Any("error", err), slog.Any("identifier", identifier.Value))
		resp.WriteHeader(500)
		return
	}
	if len(res.Headers) > 0 {
		for k, v := range res.Headers {
			resp.Header().Add(k, strings.Join(v, ",'"))
		}
	}
	resp.WriteHeader(200)
	defer (*res.Body).Close()

	buff := make([]byte, 1024)
	io.CopyBuffer(resp, *res.Body, buff)
}
