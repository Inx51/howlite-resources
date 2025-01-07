package resource

import (
	"fmt"
)

type AlreadyExistsError struct {
	Identifier *ResourceIdentifier
}

func (r AlreadyExistsError) Error() string {
	return fmt.Sprintf("Resource with identifier %s already exists", *r.Identifier.Value)
}

type NotFoundError struct {
	Identifier *ResourceIdentifier
}

func (r NotFoundError) Error() string {
	return fmt.Sprintf("Could not find resource with identifier %s", *r.Identifier.Value)
}
