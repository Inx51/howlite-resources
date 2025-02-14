package resource

import (
	"io"
)

type Resource struct {
	Identifier *ResourceIdentifier
	Headers    *map[string][]string
	Body       *io.ReadCloser
}

func NewResource(identifier *ResourceIdentifier, headers *map[string][]string, body *io.ReadCloser) *Resource {
	return &Resource{
		Identifier: identifier,
		Headers:    headers,
		Body:       body,
	}
}

func (resource *Resource) ToReadCloser() io.ReadCloser {
	return nil
}
