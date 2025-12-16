package resource

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/inx51/howlite-resources/logger"
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

func (resourceHeaders *ResourceHeaders) Add(ctx context.Context, key string, values []string) {
	if isReservedResponseHeader(key) {
		logger.Debug(ctx, "Header filtered from resource", "header", key)
		return
	}

	resourceHeaders.headers[key] = values
}

func (resourceHeaders *ResourceHeaders) Headers() *map[string][]string {
	return &resourceHeaders.headers
}

func (resourceHeaders *ResourceHeaders) LoadHeaders(reader io.ReadCloser) error {
	headers, err := getHeadersFromStream(reader)
	if err != nil {
		return err
	}

	resourceHeaders.headers = *headers
	return nil
}

func isReservedResponseHeader(headerName string) bool {
	var reservedResponseHeaders = []string{
		"content-length",
		"transfer-encoding",
		"connection",
		"upgrade",
		"trailer",
		"te",
		"host",
		"server",
		"date",
		"location",
		"retry-after",
		"vary",
		"warning",
		"status",
		"set-cookie",
		"set-cookie2",
		"www-authenticate",
		"proxy-authenticate",
		"authorization",
		"proxy-authorization",
		"age",
		"cache-control",
		"expires",
		"last-modified",
		"etag",
		"if-match",
		"if-none-match",
		"if-modified-since",
		"if-unmodified-since",
		"if-range",
		"accept-ranges",
		"content-range",
		"content-encoding",
		"content-language",
		"content-disposition",
		"content-md5",
		"accept",
		"accept-charset",
		"accept-encoding",
		"accept-language",
		"user-agent",
		"referer",
		"origin",
		"range",
		"expect",
		"max-forwards",
		"from",
		"access-control-allow-origin",
		"access-control-allow-credentials",
		"access-control-expose-headers",
		"access-control-max-age",
		"access-control-allow-methods",
		"access-control-allow-headers",
		"access-control-request-method",
		"access-control-request-headers",
		"upgrade-insecure-requests",
		"x-forwarded-for",
		"x-forwarded-host",
		"x-forwarded-proto",
		"x-real-ip",
		"via",
		"forwarded",
	}

	return slices.Contains(reservedResponseHeaders, strings.ToLower(headerName))
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

	length := binary.LittleEndian.Uint64(headerLengthBytes)
	fmt.Println(length)
	return length, nil
}
