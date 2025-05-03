package storage

import (
	"context"
	"io"

	"github.com/inx51/howlite/resources/resource"
)

type Storage interface {
	GetName() string
	RemoveResourceContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) error
	NewResourceWriterContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (io.WriteCloser, error)
	ResourceExistsContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (bool, error)
	GetResourceContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (io.ReadCloser, error)
}
