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
	logger.Debug(ctx, "Trying to read resource from file", "resourceIdentifier", resourceIdentifier.Identifier(), "file", path)
	osOpenCtx, span := tracer.StartDebugSpan(ctx, "os.open")
	reader, err := os.Open(path)
	tracer.SetDebugAttributes(osOpenCtx, span, attribute.String("path", path))
	tracer.SafeEndSpan(span)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logger.Debug(ctx, "Could not find resource file", "resourceIdentifier", resourceIdentifier.Identifier(), "file", path)
			return nil, err
		}
		logger.Error(ctx, "Failed to read resource from file", "resourceIdentifier", resourceIdentifier.Identifier(), "file", path)
		return nil, err
	}
	resource, err := resource.LoadResource(resourceIdentifier, reader)
	logger.Debug(ctx, "Read resource from file", "resourceIdentifier", resourceIdentifier.Identifier(), "file", path)
	return resource, err
}

func (fileSystem *Storage) RemoveResource(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) error {
	path := fileSystem.resourcePath(resourceIdentifier)
	logger.Debug(ctx, "Trying to remove resource file", "resourceIdentifier", resourceIdentifier.Identifier(), "file", path)
	logger.Debug(ctx, "Trying to read resource from file", "resourceIdentifier", resourceIdentifier.Identifier(), "file", path)
	osRemoveCtx, span := tracer.StartDebugSpan(ctx, "os.remove")
	err := os.Remove(path)
	tracer.SetDebugAttributes(osRemoveCtx, span, attribute.String("path", path))
	tracer.SafeEndSpan(span)
	if err != nil {
		logger.Debug(ctx, "Removing resource file failed with unhandled error", "resourceIdentifier", resourceIdentifier.Identifier(), "file", path, "error", err)
		return err
	}
	logger.Debug(ctx, "Successfully removed resource file", "resourceIdentifier", resourceIdentifier.Identifier(), "file", path)
	return nil
}

func (fileSystem *Storage) ResourceExists(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (bool, error) {
	path := fileSystem.resourcePath(resourceIdentifier)
	logger.Debug(ctx, "Trying to validate if resource file exists", "resourceIdentifier", resourceIdentifier.Identifier(), "file", path)
	osStatCtx, span := tracer.StartDebugSpan(ctx, "os.stat")
	_, err := os.Stat(path)
	tracer.SetDebugAttributes(osStatCtx, span, attribute.String("path", path))
	tracer.SafeEndSpan(span)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logger.Debug(ctx, "Could not find resource file", "resourceIdentifier", resourceIdentifier.Identifier(), "file", path)
			return false, nil
		}
		logger.Error(ctx, "Failed to find resource with unhandled error", "error", err)
		return false, err
	}
	logger.Debug(ctx, "Found resource file", "resourceIdentifier", resourceIdentifier.Identifier(), "file", path)
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
	logger.Debug(ctx, "Trying to create new writer for resource file", "resourceIdentifier", resource.Identifier.Identifier(), "file", path)
	osCreateCtx, span := tracer.StartDebugSpan(ctx, "os.create")
	writer, err := os.Create(path)
	tracer.SetDebugAttributes(osCreateCtx, span, attribute.String("path", path))
	tracer.SafeEndSpan(span)
	if err != nil {
		logger.Error(ctx, "Failed to create new writer for resource file", "resourceIdentifier", resource.Identifier.Identifier(), "path", path, "error", err)
		return err
	}
	logger.Debug(ctx, "Successfully created new writer for resource file", "resourceIdentifier", resource.Identifier.Identifier(), "file", path)
	resource.Write(writer)
	writer.Close()
	return nil
}

func (fileSystem *Storage) resourcePath(resourceIdentifier *resource.ResourceIdentifier) string {
	return fileSystem.StoragePath + "/" + resourceIdentifier.ToUniqueFilename()
}
