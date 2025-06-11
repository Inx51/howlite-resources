package storage

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/inx51/howlite/resources/config"
	"github.com/inx51/howlite/resources/storage/filesystem"
	"github.com/inx51/howlite/resources/storage/s3"
)

func Create(logger *slog.Logger, config config.StorageProvider) (Storage, error) {
	switch strings.ToLower(config.STORAGE_PROVIDER) {
	case "filesystem":
		return filesystem.NewStorage(config.STORAGE_PROVIDER_FILESYSTEM, logger), nil
	case "s3":
		return s3.NewStorage(config.STORAGE_PROVIDER_S3, logger)
	default:
		return nil, fmt.Errorf("unsupported storage provider: %s", config.STORAGE_PROVIDER)
	}
}
