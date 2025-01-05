package resource

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/inx51/howlite/resources/config"
	"github.com/inx51/howlite/resources/hash"
	"github.com/vmihailenco/msgpack/v5"
)

type Resource struct {
	Identifier *ResourceIdentifier
	Headers    map[string][]string
	Body       *io.ReadCloser
}

func New(identifier *ResourceIdentifier, headers map[string][]string, body *io.ReadCloser) Resource {
	return Resource{
		Identifier: identifier,
		Headers:    headers,
		Body:       body,
	}
}

func Get(identifier *ResourceIdentifier) (*Resource, error) {
	path := getPath(identifier)
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, NotFoundError{Identifier: identifier}
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
		resource := New(identifier, headers, &readCloser)
		return &resource, nil
	} else {
		resource := New(identifier, nil, &readCloser)
		return &resource, nil
	}
}

func Exists(identifier *ResourceIdentifier) bool {
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

func Create(res *Resource) error {
	if Exists(res.Identifier) {
		return AlreadyExistsError{Identifier: res.Identifier}
	}

	//Output
	path := getPath(res.Identifier)
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	//Headers
	res.Headers = filterHeaders(&res.Headers)
	headersLength := len(res.Headers)
	if headersLength == 0 {
		file.Write(make([]byte, 8))
	} else {
		msgPackedHeaders, err := msgpack.Marshal(res.Headers)
		if err != nil {
			panic(err)
		}
		headersMsgPackLength := len(msgPackedHeaders)
		binary.Write(file, binary.LittleEndian, uint64(headersMsgPackLength))
		file.Write(msgPackedHeaders)
	}

	//Body
	buff := make([]byte, 1024)
	readCloser := io.NopCloser(*res.Body)
	_, err = io.CopyBuffer(file, readCloser, buff)
	if err != nil {
		panic(err)
	}

	return nil
}

func filterHeaders(headers *map[string][]string) map[string][]string {
	forbiddenHeaders := []string{"host", "accept-encoding", "connection", "accepts", "user-agent"}
	var result = make(map[string][]string)
	for k, v := range *headers {
		if slices.Contains(forbiddenHeaders, strings.ToLower(k)) {
			continue
		}

		result[k] = v
	}

	return result
}

func getPath(identifier *ResourceIdentifier) string {
	hashedFileName := hash.Base64HashString(*identifier.Value)
	return config.Instance.Storage.Path + "\\" + hashedFileName + ".bin"
}
