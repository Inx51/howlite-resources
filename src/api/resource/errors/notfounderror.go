package errors

import (
	"fmt"
)

type NotFoundError struct {
	Identifier string
}

func (r NotFoundError) Error() string {
	return fmt.Sprintf("Could not find resource with identifier %s", r.Identifier)
}
