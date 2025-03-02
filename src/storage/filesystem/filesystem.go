package filesystem

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"io"
	"log/slog"
	"os"

	"github.com/inx51/howlite/resources/resource"
)

type FileSystem struct {
	logger      *slog.Logger
	StoragePath string
}

func NewStorage(storagePath string, logger *slog.Logger) *FileSystem {
	return &FileSystem{
		StoragePath: storagePath,
		logger:      logger,
	}
}

func (fileSystem *FileSystem) GetName() string {
	return "FileSystemStorage"
}

func (fileSystem *FileSystem) RemoveResource(resourceIdentifier *resource.ResourceIdentifier) error {
	path := fileSystem.resourcePath(*resourceIdentifier)
	fileSystem.logger.Debug("Trying to remove resource file", "resourceIdentifier", resourceIdentifier.Value, "file", path)
	err := os.Remove(path)
	if err != nil {
		fileSystem.logger.Debug("Removing resource file failed with unhandled error", "resourceIdentifier", resourceIdentifier.Value, "file", path, "error", err)
		return err
	}
	fileSystem.logger.Debug("Successfully removed resource file", "resourceIdentifier", resourceIdentifier.Value, "file", path)
	return nil
}

func (fileSystem *FileSystem) NewResourceWriter(resourceIdentifier *resource.ResourceIdentifier) (io.WriteCloser, error) {
	path := fileSystem.resourcePath(*resourceIdentifier)
	fileSystem.logger.Debug("Trying to create new writer for resource file", "resourceIdentifier", resourceIdentifier.Value, "file", path)
	writer, err := os.Create(path)
	if err != nil {
		fileSystem.logger.Error("Failed to create new writer for resource file", "resourceIdentifier", resourceIdentifier.Value, "path", path, "error", err)
		return nil, err
	}
	fileSystem.logger.Debug("Successfully created new writer for resource file", "resourceIdentifier", resourceIdentifier.Value, "file", path)
	return writer, nil
}

func (fileSystem *FileSystem) ResourceExists(resourceIdentifier *resource.ResourceIdentifier) (bool, error) {
	path := fileSystem.resourcePath(*resourceIdentifier)
	fileSystem.logger.Debug("Trying to validate if resource file exists", "resourceIdentifier", resourceIdentifier.Value, "file", path)
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fileSystem.logger.Debug("Could not find resource file", "resourceIdentifier", resourceIdentifier.Value, "file", path)
			return false, err
		}
		fileSystem.logger.Error("Failed to find resource with unhandled error", "error", err)
		return false, err
	}
	fileSystem.logger.Debug("Found resource file", "resourceIdentifier", resourceIdentifier.Value, "file", path)
	return true, nil
}

func (fileSystem *FileSystem) GetResource(resourceIdentifier *resource.ResourceIdentifier) (io.ReadCloser, error) {
	path := fileSystem.resourcePath(*resourceIdentifier)
	fileSystem.logger.Debug("Trying to read resource from file", "resourceIdentifier", resourceIdentifier.Value, "file", path)
	reader, err := os.Open(path)
	if err != nil {
		fileSystem.logger.Error("Failed to read resource from file", "resourceIdentifier", resourceIdentifier.Value, "file", path)
		return nil, err
	}
	fileSystem.logger.Debug("Read resource from file", "resourceIdentifier", resourceIdentifier.Value, "file", path)
	return reader, nil
}

func (fileSystem *FileSystem) resourcePath(resourceIdentifier resource.ResourceIdentifier) string {
	return fileSystem.StoragePath + "\\" + identifierToStringVersion(resourceIdentifier) + ".bin"
}

func identifierToStringVersion(resourceIdentifier resource.ResourceIdentifier) string {
	var encBytes = md5.Sum([]byte(*resourceIdentifier.Value))
	return base64.URLEncoding.EncodeToString(encBytes[:])
}
