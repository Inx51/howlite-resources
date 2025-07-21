package storage

import (
	"log/slog"
	"strings"
	"testing"

	"github.com/inx51/howlite/resources/config"
	"github.com/inx51/howlite/resources/testing/utilities/logging"
	"github.com/stretchr/testify/assert"
)

func TestCreateShouldReturnFilestorageGivenFilestorageProviderConfigured(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(&logging.TestingLogWriter{}, nil))
	storageproviderConfig := &config.StorageProvider{
		STORAGE_PROVIDER: "filesystem",
	}

	storage, _ := Create(logger, *storageproviderConfig)
	storageName := strings.ToLower(storage.GetName())

	assert.Equal(t, storageName, "filesystem")
}

func TestCreateShouldReturnS3StorageGivenS3StorageProviderConfigured(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(&logging.TestingLogWriter{}, nil))
	storageproviderConfig := &config.StorageProvider{
		STORAGE_PROVIDER: "s3",
	}

	storage, _ := Create(logger, *storageproviderConfig)
	storageName := strings.ToLower(storage.GetName())

	assert.Equal(t, storageName, "s3")
}

func TestShouldThrowErrorGivenNoStorageProviderFound(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(&logging.TestingLogWriter{}, nil))
	storageproviderConfig := &config.StorageProvider{
		STORAGE_PROVIDER: "NOT_IMPLEMENTED",
	}

	_, err := Create(logger, *storageproviderConfig)

	assert.Error(t, err)
}
