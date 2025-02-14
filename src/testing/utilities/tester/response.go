package tester

import (
	"io"
	"net/http"
	"strings"
)

type Response struct {
	StatusCode int
	Headers    Headers
	Body       []byte
}

func NewResponse(response http.Response) Response {
	headers := make(map[string]string)
	for k, v := range response.Header {
		headers[k] = strings.Join(v, ",")
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	return Response{
		StatusCode: response.StatusCode,
		Headers: Headers{
			headers: headers,
		},
		Body: bodyBytes,
	}
}
