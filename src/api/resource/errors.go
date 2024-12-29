package resource

import (
	"errors"
	"fmt"
)

var (
	ErrResourceNotFound = errNotFound()
)

type ResourceNotFound struct {
	Identifier string
}

// Error implements error.
func (r ResourceNotFound) Error() string {
	return fmt.Sprintf("Could not find resource with identifier %s", r.Identifier)
}

func errNotFound(r ResourceNotFound) error {
	return errors.New()
}
