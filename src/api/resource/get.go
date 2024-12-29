package resource

import (
	"errors"
	"io/fs"
	"os"
)

func Get(identifier string) (Resource, error) {

	file, err := os.Open("")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Resource{}, ResourceNotFound{Identifier: identifier}
		}

		return Resource{}, err
	}

	return Resource{
		Identifier: identifier,
		Body:       file,
	}, nil
}
