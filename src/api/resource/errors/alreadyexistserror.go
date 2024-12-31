package errors

import (
	"fmt"
)

type AlreadyExistsError struct {
	Identifier string
}

func (r AlreadyExistsError) Error() string {
	return fmt.Sprintf("Resource with identifier %s already exists", r.Identifier)
}
