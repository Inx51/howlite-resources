package repository

import (
	"encoding/binary"
	"io"
	"slices"
	"strings"

	"github.com/inx51/howlite/resources/resource"
	"github.com/inx51/howlite/resources/storage"
	"github.com/vmihailenco/msgpack/v5"
)

type Repository struct {
	Storage storage.Storage
}

func NewRepository(storage storage.Storage) *Repository {
	return &Repository{
		Storage: storage,
	}
}

func (repository *Repository) GetResource(resourceIdentifier *resource.ResourceIdentifier) (*resource.Resource, error) {
	resourceStream, err := repository.Storage.GetResource(resourceIdentifier)
	if err != nil {
		panic(err)
	}

	headers := buildHeaders(&resourceStream)
	resource := resource.NewResource(resourceIdentifier, &headers, &resourceStream)

	return resource, nil
}

func (repository *Repository) ResourceExists(resourceIdentifier *resource.ResourceIdentifier) (bool, error) {
	exists, err := repository.Storage.ResourceExists(resourceIdentifier)
	if err != nil {
		panic(err)
	}

	return exists, nil
}

func (repository *Repository) SaveResource(resource *resource.Resource) error {

	resourceStream, err := repository.Storage.NewResourceWriter(resource.Identifier)
	if err != nil {
		panic(err)
	}

	writeHeaders(resource, resourceStream)
	writeBody(resource, resourceStream)

	return nil
}

func (repository *Repository) RemoveResource(resourceIdentifier *resource.ResourceIdentifier) error {
	err := repository.Storage.RemoveResource(resourceIdentifier)
	if err != nil {
		panic(err)
	}

	return nil
}

func writeHeaders(resource *resource.Resource, resourceStream io.WriteCloser) {
	resource.Headers = filterHeaders(resource.Headers)

	headers := *resource.Headers
	headersLength := len(headers)
	if headersLength == 0 {
		resourceStream.Write(make([]byte, 8))
	} else {
		msgPackedHeaders, err := msgpack.Marshal(headers)
		if err != nil {
			panic(err)
		}
		headersMsgPackLength := len(msgPackedHeaders)
		binary.Write(resourceStream, binary.LittleEndian, uint64(headersMsgPackLength))
		resourceStream.Write(msgPackedHeaders)
	}
}

func writeBody(resource *resource.Resource, resourceStream io.WriteCloser) {
	buff := make([]byte, 1024)
	readCloser := io.NopCloser(*resource.Body)
	_, err := io.CopyBuffer(resourceStream, readCloser, buff)
	if err != nil {
		panic(err)
	}

	resourceStream.Close()
}

func buildHeaders(resourceStream *io.ReadCloser) map[string][]string {
	var headers map[string][]string = nil

	headerLengthBytes := make([]byte, 8)
	resStream := *resourceStream
	_, err := resStream.Read(headerLengthBytes)
	// TODO: Deal with this..
	if err != nil {
		panic(err)
	}
	headerLength := binary.LittleEndian.Uint64(headerLengthBytes)
	if headerLength > 0 {
		headerBytes := make([]byte, headerLength)
		resStream.Read(headerBytes)
		msgpack.Unmarshal(headerBytes, &headers)
	}

	return headers
}

func filterHeaders(headers *map[string][]string) *map[string][]string {
	forbiddenHeaders := []string{"host", "accept-encoding", "connection", "accepts", "user-agent"}
	var result = make(map[string][]string)
	for k, v := range *headers {
		if slices.Contains(forbiddenHeaders, strings.ToLower(k)) {
			continue
		}

		result[k] = v
	}

	return &result
}
