package fakes

import (
	"github.com/inx51/howlite/resources/resource"
	"github.com/inx51/howlite/resources/storage"
)

type FakeStorage struct {
	resources map[string]resource.Resource
}

func NewStorage() storage.Storage {
	return FakeStorage{}
}

func (s FakeStorage) Store(resource *resource.Resource) error {
	s.resources[*(resource).Identifier.Value] = *resource

	return nil
}

func (s FakeStorage) Load(identifier *resource.ResourceIdentifier) (*resource.Resource, error) {

	resource, ok := s.resources[*identifier.Value]

	if ok {
		return &resource, nil
	}

	return nil, storage.NotFoundError{
		Identifier: resource.Identifier,
	}
}
