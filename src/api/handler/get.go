package handler

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/inx51/howlite/resources/api/resource"
)

func GetResource(resp http.ResponseWriter, req *http.Request) {
	identifier := resource.GetIdentifier(req.URL.Path)
	res, err := resource.Get(identifier)
	if err != nil {
		if errors.Is(err, resource.ResourceNotFound) {
			resp.WriteHeader(404)
			slog.Warn("Failed to get resource", err)
		}
		resp.WriteHeader(500)
		slog.Error("Unhandled error", err)
	}
	defer res.Body.Close()
	buff := make([]byte, 1024)
	io.CopyBuffer(resp, res.Body, buff)
	resp.WriteHeader(201)
}
