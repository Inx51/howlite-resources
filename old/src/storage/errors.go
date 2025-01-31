package storage

import (
	"fmt"

	"github.com/inx51/howlite/resources/resource"
)

type NotFoundError struct {
	Identifier *resource.ResourceIdentifier
}

func (r NotFoundError) Error() string {
	return fmt.Sprintf("Could not find resource with identifier %s", *r.Identifier.Value)
}
