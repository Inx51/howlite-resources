package storage

import "github.com/inx51/howlite/resources/resource"

type Storage interface {
	Store(resource *resource.Resource) error
	Load(identifier *resource.ResourceIdentifier) (*resource.Resource, error)
}

func Create() Storage {
	return nil
}
