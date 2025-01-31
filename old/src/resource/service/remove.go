package service

import (
	"errors"
	"os"

	"github.com/inx51/howlite/resources/resource"
	"github.com/inx51/howlite/resources/storage"
)

func Remove(identifier *resource.ResourceIdentifier, storage *storage.Storage) error {
	path := getPath(identifier)
	err := os.Remove(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return resource.NotFoundError{Identifier: identifier}
		}
		panic(err)
	}
	return nil
}
