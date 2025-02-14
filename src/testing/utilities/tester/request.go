package tester

import (
	"io"
	"strings"
)

type Request struct {
	Path    string
	Method  string
	Headers *Headers
	Body    *RequestBody
}

type RequestBody struct {
	headers *Headers
	reader  io.Reader
}

func NewRequest() Request {
	headers := NewHeaders()
	body := RequestBody{
		headers: headers,
	}

	return Request{
		Headers: headers,
		Body:    &body,
	}
}

func (body *RequestBody) SetString(json string) {
	body.reader = strings.NewReader(json)
}
