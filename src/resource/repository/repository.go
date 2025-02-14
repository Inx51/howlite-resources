package repository

import (
	"github.com/inx51/howlite/resources/resource"
	"github.com/inx51/howlite/resources/resource/services"
	"github.com/inx51/howlite/resources/storage"
)

type Repository struct {
	Storage storage.Storage
}

func NewRepository(storage storage.Storage) *Repository {
	return &Repository{
		Storage: storage,
	}
}

func (repository *Repository) GetResource(resourceIdentifier *resource.ResourceIdentifier) (*resource.Resource, error) {
	resourceStream, err := repository.Storage.GetResource(resourceIdentifier)
	if err != nil {
		panic(err)
	}

	headers := services.GetHeadersFromStream(&resourceStream)
	resource := resource.NewResource(resourceIdentifier, headers, &resourceStream)

	return resource, nil
}

func (repository *Repository) ResourceExists(resourceIdentifier *resource.ResourceIdentifier) (bool, error) {
	exists, err := repository.Storage.ResourceExists(resourceIdentifier)
	if err != nil {
		panic(err)
	}

	return exists, nil
}

func (repository *Repository) SaveResource(resource *resource.Resource) error {

	resourceStream, err := repository.Storage.NewResourceWriter(resource.Identifier)
	if err != nil {
		panic(err)
	}

	//Write headers
	headers := services.FilterForValidResponseHeaders(resource.Headers)
	services.WriteHeaders(&resourceStream, headers)

	//Write body
	services.WriteBody(&resourceStream, resource.Body)

	return nil
}

func (repository *Repository) RemoveResource(resourceIdentifier *resource.ResourceIdentifier) error {
	err := repository.Storage.RemoveResource(resourceIdentifier)
	if err != nil {
		panic(err)
	}

	return nil
}
