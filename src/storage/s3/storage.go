package s3

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"strings"

	"github.com/inx51/howlite/resources/config"
	"github.com/inx51/howlite/resources/resource"
)

type Storage struct {
	client Client
	config config.S3Configuration
	logger *slog.Logger
}

func NewStorage(config config.S3Configuration, client Client, log *slog.Logger) (*Storage, error) {
	return &Storage{
		client: client,
		config: config,
		logger: log,
	}, nil
}

func (storage *Storage) GetName() string {
	return "S3"
}

func (storage *Storage) RemoveResourceContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) error {
	return storage.client.RemoveObjectContext(ctx, &storage.config.BUCKET, resourceIdentifier.Value)
}

func (storage *Storage) NewResourceWriterContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (io.WriteCloser, error) {
	strategy := strings.ToLower(storage.config.UPLOAD_STRATEGY)
	switch strategy {
	case "multipart":
		return NewMultipartWriterWithContext(ctx, storage.client, storage.config.BUCKET, storage.config.ACCESS_KEY, storage.config.MULTIPART_PART_UPLOAD_SIZE, &bytes.Buffer{}), nil
	case "singlepart":
		return NewSinglepartWriterWithContext(ctx, storage.client, storage.config.BUCKET, storage.config.ACCESS_KEY, &bytes.Buffer{}), nil
	default:
		return nil, NewUndefinedStrategyError(strategy)
	}
}

func (storage *Storage) ResourceExistsContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (bool, error) {
	return storage.client.ObjectExistsContext(ctx, &storage.config.BUCKET, resourceIdentifier.Value)
}

func (storage *Storage) GetResourceContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (io.ReadCloser, error) {
	return storage.client.GetObjectContext(ctx, &storage.config.BUCKET, resourceIdentifier.Value)
}
