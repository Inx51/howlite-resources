package filesystem

import (
	"context"
	"errors"
	"os"

	"github.com/inx51/howlite-resources/configuration"
	"github.com/inx51/howlite-resources/logger"
	"github.com/inx51/howlite-resources/resource"
	"github.com/inx51/howlite-resources/storage"
	"github.com/inx51/howlite-resources/tracer"
	"go.opentelemetry.io/otel/attribute"
)

type Storage struct {
	StoragePath string
}

func (fileSystem *Storage) GetResource(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (*resource.Resource, error) {
	path := fileSystem.resourcePath(resourceIdentifier)
	logger.Debug(ctx, "trying to read file", "resource.identifier", resourceIdentifier.Identifier(), "file.path", path)

	osOpenCtx, span := tracer.StartDebugSpan(ctx, "os.open")
	tracer.SetDebugAttributes(osOpenCtx, span,
		attribute.String("file.path", path),
		attribute.String("resource.identifier", resourceIdentifier.Identifier()),
	)
	reader, err := os.Open(path)
	if err != nil {
		span.RecordError(err)
	}
	tracer.SafeEndSpan(span)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logger.Debug(ctx, "file not found", "resource.identifier", resourceIdentifier.Identifier(), "file.path", path)
			return nil, err
		}
		logger.Error(ctx, "failed to open file", "resource.identifier", resourceIdentifier.Identifier(), "file.path", path, "error", err)
		return nil, err
	}

	resource, err := resource.LoadResource(resourceIdentifier, reader)
	if err != nil {
		logger.Error(ctx, "failed to load resource from file", "resource.identifier", resourceIdentifier.Identifier(), "error", err)
		return nil, err
	}
	logger.Debug(ctx, "successfully read file", "resource.identifier", resourceIdentifier.Identifier(), "file.path", path)
	return resource, err
}

func (fileSystem *Storage) RemoveResource(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) error {
	path := fileSystem.resourcePath(resourceIdentifier)
	logger.Debug(ctx, "trying to remove file", "resource.identifier", resourceIdentifier.Identifier(), "file.path", path)

	osRemoveCtx, span := tracer.StartDebugSpan(ctx, "os.remove")
	tracer.SetDebugAttributes(osRemoveCtx, span,
		attribute.String("file.path", path),
		attribute.String("resource.identifier", resourceIdentifier.Identifier()),
	)
	err := os.Remove(path)
	defer tracer.SafeEndSpan(span)

	if err != nil {
		span.RecordError(err)
		logger.Error(ctx, "failed to remove file", "resource.identifier", resourceIdentifier.Identifier(), "file.path", path, "error", err)
		return err
	}
	logger.Debug(ctx, "successfully removed file", "resource.identifier", resourceIdentifier.Identifier(), "file.path", path)
	return nil
}

func (fileSystem *Storage) ResourceExists(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (bool, error) {
	path := fileSystem.resourcePath(resourceIdentifier)
	logger.Debug(ctx, "checking if file exists", "resource.identifier", resourceIdentifier.Identifier(), "file.path", path)

	osStatCtx, span := tracer.StartDebugSpan(ctx, "os.stat")
	defer tracer.SafeEndSpan(span)
	tracer.SetDebugAttributes(osStatCtx, span,
		attribute.String("file.path", path),
		attribute.String("resource.identifier", resourceIdentifier.Identifier()),
	)
	_, err := os.Stat(path)

	if err != nil {
		span.RecordError(err)
		if errors.Is(err, os.ErrNotExist) {
			logger.Debug(ctx, "file not found", "resource.identifier", resourceIdentifier.Identifier(), "file.path", path)
			return false, nil
		}
		logger.Error(ctx, "failed to stat file", "resource.identifier", resourceIdentifier.Identifier(), "file.path", path, "error", err)
		return false, err
	}

	logger.Debug(ctx, "file exists", "resource.identifier", resourceIdentifier.Identifier(), "file.path", path)
	return true, nil
}

func NewStorage(configuration *configuration.FilesystemConfiguration) storage.Storage {
	return &Storage{StoragePath: configuration.PATH}
}

func (fileSystem *Storage) GetName() string {
	return "filesystem"
}

func (fileSystem *Storage) SaveResource(ctx context.Context, resource *resource.Resource) error {
	path := fileSystem.resourcePath(resource.Identifier)
	logger.Debug(ctx, "trying to create file", "resource.identifier", resource.Identifier.Identifier(), "file.path", path)

	osCreateCtx, span := tracer.StartDebugSpan(ctx, "os.create")
	tracer.SetDebugAttributes(osCreateCtx, span,
		attribute.String("file.path", path),
		attribute.String("resource.identifier", resource.Identifier.Identifier()),
	)
	writer, err := os.Create(path)
	if err != nil {
		span.RecordError(err)
	}
	tracer.SafeEndSpan(span)

	if err != nil {
		logger.Error(ctx, "failed to create file", "resource.identifier", resource.Identifier.Identifier(), "file.path", path, "error", err)
		return err
	}

	logger.Debug(ctx, "successfully created file", "resource.identifier", resource.Identifier.Identifier(), "file.path", path)
	resource.Write(writer)
	writer.Close()
	return nil
}

func (fileSystem *Storage) resourcePath(resourceIdentifier *resource.ResourceIdentifier) string {
	return fileSystem.StoragePath + "/" + resourceIdentifier.ToUniqueFilename()
}
