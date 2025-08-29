package resource

import (
	"encoding/binary"
	"io"

	"github.com/vmihailenco/msgpack/v5"
)

type ResourceHeaders struct {
	headers map[string][]string
}

func NewResourceHeaders() *ResourceHeaders {
	return &ResourceHeaders{
		headers: make(map[string][]string),
	}
}

func (resourceHeaders *ResourceHeaders) Add(key string, values []string) {
	resourceHeaders.headers[key] = values
}

func (resourceHeaders *ResourceHeaders) Headers() *map[string][]string {
	return &resourceHeaders.headers
}

func (resourceHeaders *ResourceHeaders) writeHeaders(writer io.Writer) error {
	headersLength := len(resourceHeaders.headers)
	if headersLength == 0 {
		writer.Write(make([]byte, 8))
	} else {
		msgPackedHeaders, err := msgpack.Marshal(resourceHeaders.headers)
		if err != nil {
			return err
		}
		headersMsgPackLength := len(msgPackedHeaders)
		binary.Write(writer, binary.LittleEndian, uint64(headersMsgPackLength))
		writer.Write(msgPackedHeaders)
	}
	return nil
}

func (resourceHeaders *ResourceHeaders) LoadHeaders(reader io.ReadCloser) error {
	headers, err := getHeadersFromStream(reader)
	if err != nil {
		return err
	}

	resourceHeaders.headers = *headers
	return nil
}

func getHeadersFromStream(stream io.ReadCloser) (*map[string][]string, error) {
	var headers map[string][]string
	headersLength, err := getHeadersLengthFromStream(stream)
	if err != nil {
		return nil, err
	}
	if headersLength == 0 {
		return &headers, nil
	}
	headersBytes := make([]byte, headersLength)
	_, err = stream.Read(headersBytes)
	if err != nil {
		return nil, err
	}

	err = msgpack.Unmarshal(headersBytes, &headers)
	if err != nil {
		return nil, err
	}

	return &headers, nil
}

func getHeadersLengthFromStream(stream io.ReadCloser) (uint64, error) {
	headerLengthBytes := make([]byte, 8)
	_, err := stream.Read(headerLengthBytes)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint64(headerLengthBytes), nil
}
