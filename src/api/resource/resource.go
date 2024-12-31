package resource

import (
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/inx51/howlite/resources/api/config"
	"github.com/inx51/howlite/resources/api/hash"
	reserrors "github.com/inx51/howlite/resources/api/resource/errors"
)

type Resource struct {
	Body       *os.File
	Identifier *string
}

func GetIdentifier(path *string) string {
	return hash.Base64HashString(*path)
}

func Get(identifier *string) (Resource, error) {
	path := getPath(identifier)
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Resource{}, reserrors.NotFoundError{Identifier: *identifier}
		}
		return Resource{}, err
	}
	return Resource{
		Identifier: identifier,
		Body:       file,
	}, nil
}

func Exists(identifier *string) bool {
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

func Create(identifier *string, body *io.ReadCloser, headers *http.Header) error {
	if Exists(identifier) {
		return reserrors.AlreadyExistsError{Identifier: *identifier}
	}

	path := getPath(identifier)

	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buff := make([]byte, 1024)
	_, err = io.CopyBuffer(file, *body, buff)
	if err != nil {
		panic(err)
	}

	return nil
}

func getPath(identifier *string) string {
	return config.Instance.Storage.Path + "\\" + *identifier + ".bin"
}
