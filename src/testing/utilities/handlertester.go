package utilities

import (
	"net/http"

	"github.com/inx51/howlite/resources/resource/repository"
	"github.com/inx51/howlite/resources/storage"
	"github.com/inx51/howlite/resources/testing/utilities/tester"
)

func CreateHandlerParameters(tester *tester.Tester, storage storage.Storage) (http.ResponseWriter, *http.Request, *repository.Repository) {
	resp, req := tester.Build()
	repo := createRepository(storage)

	return resp, req, repo
}

func createRepository(storage storage.Storage) *repository.Repository {
	return repository.NewRepository(storage)
}
