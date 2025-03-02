package repository

import (
	"log/slog"

	"github.com/inx51/howlite/resources/resource"
	"github.com/inx51/howlite/resources/resource/services"
	"github.com/inx51/howlite/resources/storage"
)

type Repository struct {
	Storage storage.Storage
	logger  *slog.Logger
}

func NewRepository(storage *storage.Storage, logger *slog.Logger) *Repository {
	return &Repository{
		Storage: *storage,
		logger:  logger,
	}
}

func (repository *Repository) GetResource(resourceIdentifier *resource.ResourceIdentifier) (*resource.Resource, error) {
	repository.logger.Debug("Trying to get resource", "resourceIdentifier", resourceIdentifier.Value)
	resourceStream, err := repository.Storage.GetResource(resourceIdentifier)
	if err != nil {
		repository.logger.Error("Failed to get resource", "resourceIdentifier", resourceIdentifier.Value, "error", err)
		return nil, err
	}

	headers, err := services.GetHeadersFromStream(&resourceStream)
	if err != nil {
		repository.logger.Error("Failed to get resource headers", "resourceIdentifier", resourceIdentifier.Value, "error", err)
		return nil, err
	}
	resource := resource.NewResource(resourceIdentifier, headers, &resourceStream)
	repository.logger.Debug("Successfully got resource", "resourceIdentifier", resourceIdentifier.Value)
	return resource, nil
}

func (repository *Repository) ResourceExists(resourceIdentifier *resource.ResourceIdentifier) (bool, error) {
	repository.logger.Debug("Trying to validate if resource exists", "resourceIdentifier", resourceIdentifier.Value)
	exists, err := repository.Storage.ResourceExists(resourceIdentifier)
	if err != nil {
		repository.logger.Error("Failed to get resource", "resourceIdentifier", resourceIdentifier.Value, "error", err)
		return false, err
	}
	repository.logger.Debug("Successfully validated if resource exists", "resourceIdentifier", resourceIdentifier.Value, "exists", exists)
	return exists, nil
}

func (repository *Repository) SaveResource(resource *resource.Resource) error {
	repository.logger.Debug("Trying to save resource", "resourceIdentifier", resource.Identifier.Value)
	resourceStream, err := repository.Storage.NewResourceWriter(resource.Identifier)
	if err != nil {
		repository.logger.Error("Failed to save resource", "resourceIdentifier", resource.Identifier.Value, "error", err)
		return err
	}
	//Write headers
	headers := services.FilterForValidResponseHeaders(resource.Headers, repository.logger)
	err = services.WriteHeaders(&resourceStream, headers)
	if err != nil {
		repository.logger.Error("Failed to write headers", "resourceIdentifier", resource.Identifier.Value, "error", err)
		return err
	}

	//Write body
	err = services.WriteBody(&resourceStream, resource.Body)
	if err != nil {
		repository.logger.Error("Failed to write resource body", "resourceIdentifier", resource.Identifier.Value, "error", err)
		return err
	}

	repository.logger.Info("Saved resource", "resourceIdentifier", resource.Identifier.Value)

	return nil
}

func (repository *Repository) RemoveResource(resourceIdentifier *resource.ResourceIdentifier) error {
	repository.logger.Debug("Trying to remove resource", "resourceIdentifier", resourceIdentifier.Value)
	err := repository.Storage.RemoveResource(resourceIdentifier)
	if err != nil {
		repository.logger.Error("Failed to remove resource", "resourceIdentifier", resourceIdentifier.Value, "error", err)
		return err
	}
	repository.logger.Info("Removed resource", "resourceIdentifier", resourceIdentifier.Value)
	return nil
}
