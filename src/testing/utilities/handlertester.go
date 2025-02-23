package utilities

import (
	"log/slog"
	"net/http"

	"github.com/inx51/howlite/resources/resource/repository"
	"github.com/inx51/howlite/resources/storage"
	"github.com/inx51/howlite/resources/testing/utilities/tester"
)

func CreateHandlerParameters(tester *tester.Tester, storage storage.Storage) (http.ResponseWriter, *http.Request, *repository.Repository, *slog.Logger) {
	resp, req := tester.Build()
	logger := slog.New(slog.NewTextHandler(nil, nil))
	repo := repository.NewRepository(&storage, logger)

	return resp, req, repo, logger
}
