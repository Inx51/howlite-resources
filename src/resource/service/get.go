package service

import (
	"encoding/binary"
	"errors"
	"io"
	"os"

	"github.com/inx51/howlite/resources/resource"
	"github.com/inx51/howlite/resources/storage"
	"github.com/vmihailenco/msgpack/v5"
)

func Get(identifier *resource.ResourceIdentifier, storage *storage.Storage) (*resource.Resource, error) {
	path := getPath(identifier)
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, resource.NotFoundError{Identifier: identifier}
		}
		return nil, err
	}
	readCloser := io.NopCloser(file)

	headerLengthBytes := make([]byte, 8)
	readCloser.Read(headerLengthBytes)
	if err != nil {
		panic(err)
	}
	headerLength := binary.LittleEndian.Uint64(headerLengthBytes)
	if headerLength > 0 {
		headerBytes := make([]byte, headerLength)
		readCloser.Read(headerBytes)
		var headers map[string][]string
		msgpack.Unmarshal(headerBytes, &headers)
		resource := resource.New(identifier, headers, &readCloser)
		return &resource, nil
	} else {
		resource := resource.New(identifier, nil, &readCloser)
		return &resource, nil
	}
}
