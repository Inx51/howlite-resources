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

// func BuildFromStorage(identifier *ResourceIdentifier) (*Resource, error) {
// 	return nil, nil
// }

// func ExistsInStorage(identifier *ResourceIdentifier) bool {
// 	return false
// }

// func (resource *Resource) Save() {

// }

// func (resource *Resource) Load(identifier *ResourceIdentifier) (Resource, error) {

// }
