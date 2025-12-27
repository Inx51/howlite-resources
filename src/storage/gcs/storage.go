package gcs

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"github.com/inx51/howlite-resources/configuration"
	"github.com/inx51/howlite-resources/logger"
	"github.com/inx51/howlite-resources/resource"
	storageInterface "github.com/inx51/howlite-resources/storage"
	"github.com/inx51/howlite-resources/tracer"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/api/option"
)

type Storage struct {
	client        *storage.Client
	configuration configuration.GCSConfiguration
}

func (gcsStorage *Storage) GetResource(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (*resource.Resource, error) {
	objectName := resourceIdentifier.ToUniqueFilename()
	logger.Debug(ctx, "trying to download gcs object", "resource.identifier", resourceIdentifier.Identifier(), "gcs.object", objectName)

	bucket := gcsStorage.client.Bucket(gcsStorage.configuration.BUCKET)
	object := bucket.Object(objectName)

	gcsCtx, span := tracer.StartDebugSpan(ctx, "gcs.object.download")
	tracer.SetDebugAttributes(gcsCtx, span,
		attribute.String("gcs.bucket", gcsStorage.configuration.BUCKET),
		attribute.String("gcs.object", objectName),
		attribute.String("resource.identifier", resourceIdentifier.Identifier()),
	)

	reader, err := object.NewReader(gcsCtx)
	if err != nil {
		span.RecordError(err)
		tracer.SafeEndSpan(span)

		if err == storage.ErrObjectNotExist {
			logger.Debug(ctx, "gcs object not found", "resource.identifier", resourceIdentifier.Identifier(), "gcs.object", objectName)
			return nil, err
		}

		logger.Error(ctx, "failed to download gcs object", "resource.identifier", resourceIdentifier.Identifier(), "gcs.object", objectName, "error", err)
		return nil, err
	}
	defer reader.Close()

	tracer.SafeEndSpan(span)

	resource, err := resource.LoadResource(resourceIdentifier, reader)
	if err != nil {
		logger.Error(ctx, "failed to load resource from gcs object", "resource.identifier", resourceIdentifier.Identifier(), "error", err)
		return nil, err
	}

	logger.Debug(ctx, "successfully downloaded gcs object", "resource.identifier", resourceIdentifier.Identifier(), "gcs.object", objectName)
	return resource, nil
}

func (gcsStorage *Storage) RemoveResource(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) error {
	objectName := resourceIdentifier.ToUniqueFilename()
	logger.Debug(ctx, "trying to delete gcs object", "resource.identifier", resourceIdentifier.Identifier(), "gcs.object", objectName)

	bucket := gcsStorage.client.Bucket(gcsStorage.configuration.BUCKET)
	object := bucket.Object(objectName)

	gcsCtx, span := tracer.StartDebugSpan(ctx, "gcs.object.delete")
	tracer.SetDebugAttributes(gcsCtx, span,
		attribute.String("gcs.bucket", gcsStorage.configuration.BUCKET),
		attribute.String("gcs.object", objectName),
		attribute.String("resource.identifier", resourceIdentifier.Identifier()),
	)

	err := object.Delete(gcsCtx)
	if err != nil {
		span.RecordError(err)
		tracer.SafeEndSpan(span)

		if err == storage.ErrObjectNotExist {
			logger.Debug(ctx, "gcs object not found", "resource.identifier", resourceIdentifier.Identifier(), "gcs.object", objectName)
			return err
		}

		logger.Error(ctx, "failed to delete gcs object", "resource.identifier", resourceIdentifier.Identifier(), "gcs.object", objectName, "error", err)
		return err
	}

	tracer.SafeEndSpan(span)

	logger.Debug(ctx, "successfully deleted gcs object", "resource.identifier", resourceIdentifier.Identifier(), "gcs.object", objectName)
	return nil
}

func (gcsStorage *Storage) ResourceExists(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (bool, error) {
	objectName := resourceIdentifier.ToUniqueFilename()
	logger.Debug(ctx, "checking if gcs object exists", "resource.identifier", resourceIdentifier.Identifier(), "gcs.object", objectName)

	bucket := gcsStorage.client.Bucket(gcsStorage.configuration.BUCKET)
	object := bucket.Object(objectName)

	gcsCtx, span := tracer.StartDebugSpan(ctx, "gcs.object.attrs")
	tracer.SetDebugAttributes(gcsCtx, span,
		attribute.String("gcs.bucket", gcsStorage.configuration.BUCKET),
		attribute.String("gcs.object", objectName),
		attribute.String("resource.identifier", resourceIdentifier.Identifier()),
	)

	_, err := object.Attrs(gcsCtx)
	if err != nil {
		span.RecordError(err)
		tracer.SafeEndSpan(span)

		if err == storage.ErrObjectNotExist {
			logger.Debug(ctx, "gcs object not found", "resource.identifier", resourceIdentifier.Identifier(), "gcs.object", objectName)
			return false, nil
		}

		logger.Error(ctx, "failed to check gcs object existence", "resource.identifier", resourceIdentifier.Identifier(), "gcs.object", objectName, "error", err)
		return false, err
	}

	tracer.SafeEndSpan(span)

	logger.Debug(ctx, "gcs object exists", "resource.identifier", resourceIdentifier.Identifier(), "gcs.object", objectName)
	return true, nil
}

func NewStorage(config *configuration.GCSConfiguration) storageInterface.Storage {
	ctx := context.Background()

	var clientOptions []option.ClientOption

	if config.CREDENTIALS_FILE != "" {
		clientOptions = append(clientOptions, option.WithCredentialsFile(config.CREDENTIALS_FILE))
	}

	if config.PROJECT_ID != "" {
		clientOptions = append(clientOptions, option.WithQuotaProject(config.PROJECT_ID))
	}

	client, err := storage.NewClient(ctx, clientOptions...)
	if err != nil {
		panic(fmt.Sprintf("failed to create GCS client: %v", err))
	}

	return &Storage{
		client:        client,
		configuration: *config,
	}
}

func (gcsStorage *Storage) GetName() string {
	return "gcs"
}

func (gcsStorage *Storage) Close() error {
	if gcsStorage.client != nil {
		return gcsStorage.client.Close()
	}
	return nil
}

func (gcsStorage *Storage) SaveResource(ctx context.Context, resource *resource.Resource) error {
	objectName := resource.Identifier.ToUniqueFilename()
	logger.Debug(ctx, "trying to upload gcs object", "resource.identifier", resource.Identifier.Identifier(), "gcs.object", objectName)

	bucket := gcsStorage.client.Bucket(gcsStorage.configuration.BUCKET)
	object := bucket.Object(objectName)

	gcsCtx, span := tracer.StartDebugSpan(ctx, "gcs.object.upload")
	tracer.SetDebugAttributes(gcsCtx, span,
		attribute.String("gcs.bucket", gcsStorage.configuration.BUCKET),
		attribute.String("gcs.object", objectName),
		attribute.String("resource.identifier", resource.Identifier.Identifier()),
		attribute.Int("gcs.chunk_size", gcsStorage.configuration.CHUNK_SIZE),
	)

	writer := object.NewWriter(gcsCtx)
	writer.ChunkSize = gcsStorage.configuration.CHUNK_SIZE

	reader := gcsStorage.createResourceReader(ctx, resource)
	_, err := io.Copy(writer, reader)
	if err != nil {
		span.RecordError(err)
		tracer.SafeEndSpan(span)
		logger.Error(ctx, "failed to write resource to gcs object", "resource.identifier", resource.Identifier.Identifier(), "gcs.object", objectName, "error", err)
		return err
	}

	err = writer.Close()
	if err != nil {
		span.RecordError(err)
		tracer.SafeEndSpan(span)
		logger.Error(ctx, "failed to close gcs object writer", "resource.identifier", resource.Identifier.Identifier(), "gcs.object", objectName, "error", err)
		return err
	}

	tracer.SafeEndSpan(span)

	logger.Debug(ctx, "successfully uploaded gcs object", "resource.identifier", resource.Identifier.Identifier(), "gcs.object", objectName)
	return nil
}

func (gcsStorage *Storage) createResourceReader(ctx context.Context, resource *resource.Resource) io.Reader {
	pipeReader, pipeWriter := io.Pipe()
	go func() {
		defer pipeWriter.Close()
		err := resource.Write(pipeWriter)
		if err != nil {
			logger.Error(ctx, "failed to write resource to pipe", "error", err)
			pipeWriter.CloseWithError(err)
		}
	}()
	return pipeReader
}
