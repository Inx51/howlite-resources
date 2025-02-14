package filesystem

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"io"
	"os"

	"github.com/inx51/howlite/resources/resource"
)

type FileSystem struct {
	StoragePath string
}

func NewStorage(storagePath string) *FileSystem {
	return &FileSystem{
		StoragePath: storagePath,
	}
}

func (fileSystem *FileSystem) RemoveResource(resourceIdentifier *resource.ResourceIdentifier) error {
	path := fileSystem.resourcePath(*resourceIdentifier)
	return os.Remove(path)
}

func (fileSystem *FileSystem) NewResourceWriter(resourceIdentifier *resource.ResourceIdentifier) (io.WriteCloser, error) {
	path := fileSystem.resourcePath(*resourceIdentifier)
	return os.Create(path)
}

func (fileSystem *FileSystem) ResourceExists(resourceIdentifier *resource.ResourceIdentifier) (bool, error) {
	path := fileSystem.resourcePath(*resourceIdentifier)
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (fileSystem *FileSystem) GetResource(resourceIdentifier *resource.ResourceIdentifier) (io.ReadCloser, error) {
	path := fileSystem.resourcePath(*resourceIdentifier)
	return os.Open(path)
}

func (fileSystem *FileSystem) resourcePath(resourceIdentifier resource.ResourceIdentifier) string {
	return fileSystem.StoragePath + "\\" + identifierToStringVersion(resourceIdentifier) + ".bin"
}

func identifierToStringVersion(resourceIdentifier resource.ResourceIdentifier) string {
	var encBytes = md5.Sum([]byte(*resourceIdentifier.Value))
	return base64.URLEncoding.EncodeToString(encBytes[:])
}
