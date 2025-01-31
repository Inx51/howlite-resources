package storage

import (
	"io"

	"github.com/inx51/howlite/resources/resource"
)

type Storage interface {
	RemoveResource(resourceIdentifier *resource.ResourceIdentifier) error
	NewResourceWriter(resourceIdentifier *resource.ResourceIdentifier) (io.WriteCloser, error)
	ResourceExists(resourceIdentifier *resource.ResourceIdentifier) (bool, error)
	GetResource(resourceIdentifier *resource.ResourceIdentifier) (io.ReadCloser, error)
}
