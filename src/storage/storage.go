package storage

import (
	"context"
	"io"

	"github.com/inx51/howlite/resources/resource"
)

type Storage interface {
	GetName() string
	RemoveResourceWithContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) error
	NewResourceWriterWithContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (io.WriteCloser, error)
	ResourceExistsWithContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (bool, error)
	GetResourceWithContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (io.ReadCloser, error)
}
