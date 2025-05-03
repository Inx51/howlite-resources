package services

import (
	"context"
	"encoding/binary"
	"io"
	"log/slog"
	"slices"
	"strings"

	"github.com/vmihailenco/msgpack/v5"
)

func FilterForValidResponseHeadersContext(ctx context.Context, headers *map[string][]string, logger *slog.Logger) *map[string][]string {
	forbiddenHeaders := []string{"host", "accept-encoding", "connection", "accepts", "user-agent", "authorization"}
	var result = make(map[string][]string)
	for k, v := range *headers {
		if slices.Contains(forbiddenHeaders, strings.ToLower(k)) {
			logger.Debug("Header filtered", "header", k)
			continue
		}

		result[k] = v
	}

	return &result
}

func WriteHeaders(streamRef *io.WriteCloser, headers *map[string][]string) error {
	stream := *streamRef
	headersLength := len(*headers)
	if headersLength == 0 {
		stream.Write(make([]byte, 8))
	} else {
		msgPackedHeaders, err := msgpack.Marshal(headers)
		if err != nil {
			return err
		}
		headersMsgPackLength := len(msgPackedHeaders)
		binary.Write(stream, binary.LittleEndian, uint64(headersMsgPackLength))
		stream.Write(msgPackedHeaders)
	}
	return nil
}

func GetHeadersFromStream(stream *io.ReadCloser) (*map[string][]string, error) {
	var headers map[string][]string
	headersLength, err := getHeadersLengthFromStream(stream)
	if err != nil {
		return nil, err
	}
	if headersLength == 0 {
		return &headers, nil
	}
	headersBytes := make([]byte, headersLength)
	_, err = (*stream).Read(headersBytes)
	if err != nil {
		return nil, err
	}

	err = msgpack.Unmarshal(headersBytes, &headers)
	if err != nil {
		return nil, err
	}

	return &headers, nil
}

func getHeadersLengthFromStream(stream *io.ReadCloser) (uint64, error) {
	headerLengthBytes := make([]byte, 8)
	_, err := (*stream).Read(headerLengthBytes)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint64(headerLengthBytes), nil
}
