package response

import (
	"io"
	"net/http"
	"strings"
)

func WriteHeaders(headers *map[string][]string, resp http.ResponseWriter) {
	headersToWrite := *headers
	for k, values := range headersToWrite {
		if len(values) == 1 {
			resp.Header().Set(k, values[0])
		} else {
			resp.Header().Set(k, strings.Join(values, ","))
		}
	}
}

func WriteBody(body io.ReadCloser, resp http.ResponseWriter) error {
	if body == nil {
		return nil
	}
	defer body.Close()

	_, err := io.Copy(resp, body)
	if err != nil {
		return err
	}

	return nil
}
