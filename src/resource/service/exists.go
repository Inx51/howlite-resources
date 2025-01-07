package service

import (
	"errors"
	"os"

	"github.com/inx51/howlite/resources/resource"
	"github.com/inx51/howlite/resources/storage"
)

func Exists(identifier *resource.ResourceIdentifier, storage *storage.Storage) bool {
	path := getPath(identifier)

	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		}
		panic(err)
	}
	defer file.Close()
	return true
}
