package filesystem

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"io"
	"log/slog"
	"os"

	"github.com/inx51/howlite/resources/config"
	"github.com/inx51/howlite/resources/resource"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Storage struct {
	logger      *slog.Logger
	StoragePath string
	tracer      trace.Tracer
}

func NewStorage(config config.FilesystemConfiguration, logger *slog.Logger) *Storage {
	tracer := otel.Tracer("FileSystemStorage")
	return &Storage{
		StoragePath: config.PATH,
		logger:      logger,
		tracer:      tracer,
	}
}

func (fileSystem *Storage) GetName() string {
	return "FileSystemStorage"
}

func (fileSystem *Storage) RemoveResourceContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) error {
	path := fileSystem.resourcePath(*resourceIdentifier)
	fileSystem.logger.DebugContext(ctx, "Trying to remove resource file", "resourceIdentifier", resourceIdentifier.Value, "file", path)
	err := os.Remove(path)
	if err != nil {
		fileSystem.logger.DebugContext(ctx, "Removing resource file failed with unhandled error", "resourceIdentifier", resourceIdentifier.Value, "file", path, "error", err)
		return err
	}
	fileSystem.logger.InfoContext(ctx, "Successfully removed resource file", "resourceIdentifier", resourceIdentifier.Value, "file", path)
	return nil
}

func (fileSystem *Storage) NewResourceWriterContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (io.WriteCloser, error) {
	path := fileSystem.resourcePath(*resourceIdentifier)
	fileSystem.logger.DebugContext(ctx, "Trying to create new writer for resource file", "resourceIdentifier", resourceIdentifier.Value, "file", path)
	writer, err := os.Create(path)
	if err != nil {
		fileSystem.logger.ErrorContext(ctx, "Failed to create new writer for resource file", "resourceIdentifier", resourceIdentifier.Value, "path", path, "error", err)
		return nil, err
	}
	fileSystem.logger.DebugContext(ctx, "Successfully created new writer for resource file", "resourceIdentifier", resourceIdentifier.Value, "file", path)
	return writer, nil
}

func (fileSystem *Storage) ResourceExistsContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (bool, error) {
	path := fileSystem.resourcePath(*resourceIdentifier)
	fileSystem.logger.DebugContext(ctx, "Trying to validate if resource file exists", "resourceIdentifier", resourceIdentifier.Value, "file", path)
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fileSystem.logger.DebugContext(ctx, "Could not find resource file", "resourceIdentifier", resourceIdentifier.Value, "file", path)
			return false, nil
		}
		fileSystem.logger.ErrorContext(ctx, "Failed to find resource with unhandled error", "error", err)
		return false, err
	}
	fileSystem.logger.DebugContext(ctx, "Found resource file", "resourceIdentifier", resourceIdentifier.Value, "file", path)
	return true, nil
}

func (fileSystem *Storage) GetResourceContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (io.ReadCloser, error) {
	path := fileSystem.resourcePath(*resourceIdentifier)
	fileSystem.logger.DebugContext(ctx, "Trying to read resource from file", "resourceIdentifier", resourceIdentifier.Value, "file", path)
	reader, err := os.Open(path)
	if err != nil {
		fileSystem.logger.ErrorContext(ctx, "Failed to read resource from file", "resourceIdentifier", resourceIdentifier.Value, "file", path)
		return nil, err
	}
	fileSystem.logger.DebugContext(ctx, "Read resource from file", "resourceIdentifier", resourceIdentifier.Value, "file", path)
	return reader, nil
}

func (fileSystem *Storage) resourcePath(resourceIdentifier resource.ResourceIdentifier) string {
	return fileSystem.StoragePath + "\\" + identifierToStringVersion(resourceIdentifier) + ".bin"
}

func identifierToStringVersion(resourceIdentifier resource.ResourceIdentifier) string {
	var encBytes = md5.Sum([]byte(*resourceIdentifier.Value))
	return base64.URLEncoding.EncodeToString(encBytes[:])
}
