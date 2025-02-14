package services

import (
	"encoding/binary"
	"io"
	"slices"
	"strings"

	"github.com/vmihailenco/msgpack/v5"
)

func FilterForValidResponseHeaders(headers *map[string][]string) *map[string][]string {
	forbiddenHeaders := []string{"host", "accept-encoding", "connection", "accepts", "user-agent", "authorization"}
	var result = make(map[string][]string)
	for k, v := range *headers {
		if slices.Contains(forbiddenHeaders, strings.ToLower(k)) {
			continue
		}

		result[k] = v
	}

	return &result
}

func WriteHeaders(streamRef *io.WriteCloser, headers *map[string][]string) {
	stream := *streamRef
	headersLength := len(*headers)
	if headersLength == 0 {
		stream.Write(make([]byte, 8))
	} else {
		msgPackedHeaders, err := msgpack.Marshal(headers)
		if err != nil {
			panic(err)
		}
		headersMsgPackLength := len(msgPackedHeaders)
		binary.Write(stream, binary.LittleEndian, uint64(headersMsgPackLength))
		stream.Write(msgPackedHeaders)
	}
}

func GetHeadersFromStream(stream *io.ReadCloser) *map[string][]string {
	var headers map[string][]string
	headersLength := getHeadersLengthFromStream(stream)
	if headersLength == 0 {
		return &headers
	}
	headersBytes := make([]byte, headersLength)
	_, err := (*stream).Read(headersBytes)
	if err != nil {
		panic(err)
	}

	err = msgpack.Unmarshal(headersBytes, &headers)
	if err != nil {
		panic(err)
	}

	return &headers
}

func getHeadersLengthFromStream(stream *io.ReadCloser) uint64 {
	headerLengthBytes := make([]byte, 8)
	_, err := (*stream).Read(headerLengthBytes)
	if err != nil {
		panic(err)
	}

	return binary.LittleEndian.Uint64(headerLengthBytes)
}
