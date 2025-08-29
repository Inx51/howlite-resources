package storage

import (
	"context"

	"github.com/inx51/howlite-resources/resource"
)

type Storage interface {
	RemoveResource(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) error
	SaveResource(ctx context.Context, resource *resource.Resource) error
	ResourceExists(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (bool, error)
	GetResource(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (*resource.Resource, error)
	GetName() string
}
