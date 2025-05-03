package repository

import (
	"context"
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

func (repository *Repository) GetResourceContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (*resource.Resource, error) {
	repository.logger.DebugContext(ctx, "Trying to get resource", "resourceIdentifier", resourceIdentifier.Value)
	resourceStream, err := repository.Storage.GetResourceContext(ctx, resourceIdentifier)
	if err != nil {
		repository.logger.ErrorContext(ctx, "Failed to get resource", "resourceIdentifier", resourceIdentifier.Value, "error", err)
		return nil, err
	}

	headers, err := services.GetHeadersFromStream(&resourceStream)
	if err != nil {
		repository.logger.ErrorContext(ctx, "Failed to get resource headers", "resourceIdentifier", resourceIdentifier.Value, "error", err)
		return nil, err
	}
	resource := resource.NewResource(resourceIdentifier, headers, &resourceStream)
	repository.logger.DebugContext(ctx, "Successfully got resource", "resourceIdentifier", resourceIdentifier.Value)
	return resource, nil
}

func (repository *Repository) ResourceExistsContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (bool, error) {
	repository.logger.DebugContext(ctx, "Trying to validate if resource exists", "resourceIdentifier", resourceIdentifier.Value)
	exists, err := repository.Storage.ResourceExistsContext(ctx, resourceIdentifier)
	if err != nil {
		repository.logger.ErrorContext(ctx, "Failed to get resource", "resourceIdentifier", resourceIdentifier.Value, "error", err)
		return false, err
	}
	repository.logger.DebugContext(ctx, "Successfully validated if resource exists", "resourceIdentifier", resourceIdentifier.Value, "exists", exists)
	return exists, nil
}

func (repository *Repository) SaveResourceContext(ctx context.Context, resource *resource.Resource) error {
	repository.logger.DebugContext(ctx, "Trying to save resource", "resourceIdentifier", resource.Identifier.Value)
	resourceStream, err := repository.Storage.NewResourceWriterContext(ctx, resource.Identifier)
	if err != nil {
		repository.logger.ErrorContext(ctx, "Failed to save resource", "resourceIdentifier", resource.Identifier.Value, "error", err)
		return err
	}
	//Write headers
	headers := services.FilterForValidResponseHeadersContext(ctx, resource.Headers, repository.logger)
	err = services.WriteHeaders(&resourceStream, headers)
	if err != nil {
		repository.logger.ErrorContext(ctx, "Failed to write headers", "resourceIdentifier", resource.Identifier.Value, "error", err)
		return err
	}

	//Write body
	err = services.WriteBody(&resourceStream, resource.Body)
	if err != nil {
		repository.logger.ErrorContext(ctx, "Failed to write resource body", "resourceIdentifier", resource.Identifier.Value, "error", err)
		return err
	}

	repository.logger.InfoContext(ctx, "Saved resource", "resourceIdentifier", resource.Identifier.Value)

	return nil
}

func (repository *Repository) RemoveResourceContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) error {
	repository.logger.DebugContext(ctx, "Trying to remove resource", "resourceIdentifier", resourceIdentifier.Value)
	err := repository.Storage.RemoveResourceContext(ctx, resourceIdentifier)
	if err != nil {
		repository.logger.ErrorContext(ctx, "Failed to remove resource", "resourceIdentifier", resourceIdentifier.Value, "error", err)
		return err
	}
	repository.logger.InfoContext(ctx, "Removed resource", "resourceIdentifier", resourceIdentifier.Value)
	return nil
}
