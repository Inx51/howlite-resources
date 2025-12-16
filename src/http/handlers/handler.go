package handlers

import (
	"context"
	"net/http"
)

type Handler interface {
	Method() string
	Path() string
	Handle(ctx context.Context, request *http.Request, response http.ResponseWriter) (int, error)
}
