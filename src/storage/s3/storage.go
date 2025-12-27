package s3

import (
	"context"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/inx51/howlite-resources/configuration"
	"github.com/inx51/howlite-resources/logger"
	"github.com/inx51/howlite-resources/resource"
	"github.com/inx51/howlite-resources/storage"
	"github.com/inx51/howlite-resources/tracer"
	"go.opentelemetry.io/otel/attribute"
)

type Storage struct {
	client        *s3.Client
	uploader      *manager.Uploader
	downloader    *manager.Downloader
	configuration configuration.S3Configuration
}

func (s3Storage *Storage) GetResource(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (*resource.Resource, error) {
	objectKey := resourceIdentifier.ToUniqueFilename()
	logger.Debug(ctx, "trying to download s3 object", "resource.identifier", resourceIdentifier.Identifier(), "s3.key", objectKey)

	s3Ctx, span := tracer.StartDebugSpan(ctx, "s3.get_object")
	tracer.SetDebugAttributes(s3Ctx, span,
		attribute.String("s3.bucket", s3Storage.configuration.BUCKET),
		attribute.String("s3.key", objectKey),
		attribute.String("resource.identifier", resourceIdentifier.Identifier()),
	)
	result, err := s3Storage.client.GetObject(s3Ctx, &s3.GetObjectInput{
		Bucket: aws.String(s3Storage.configuration.BUCKET),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		span.RecordError(err)
	}
	tracer.SafeEndSpan(span)

	if err != nil {
		var notFound *types.NoSuchKey
		if errors.As(err, &notFound) {
			logger.Debug(ctx, "s3 object not found", "resource.identifier", resourceIdentifier.Identifier(), "s3.key", objectKey)
			return nil, err
		}
		logger.Error(ctx, "failed to download s3 object", "resource.identifier", resourceIdentifier.Identifier(), "s3.key", objectKey, "error", err)
		return nil, err
	}

	resource, err := resource.LoadResource(resourceIdentifier, result.Body)
	if err != nil {
		result.Body.Close()
		logger.Error(ctx, "failed to load resource from s3 object", "resource.identifier", resourceIdentifier.Identifier(), "error", err)
		return nil, err
	}

	logger.Debug(ctx, "successfully downloaded s3 object", "resource.identifier", resourceIdentifier.Identifier(), "s3.key", objectKey)
	return resource, err
}

func (s3Storage *Storage) RemoveResource(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) error {
	objectKey := resourceIdentifier.ToUniqueFilename()
	logger.Debug(ctx, "trying to delete s3 object", "resource.identifier", resourceIdentifier.Identifier(), "s3.key", objectKey)

	s3Ctx, span := tracer.StartDebugSpan(ctx, "s3.delete_object")
	tracer.SetDebugAttributes(s3Ctx, span,
		attribute.String("s3.bucket", s3Storage.configuration.BUCKET),
		attribute.String("s3.key", objectKey),
		attribute.String("resource.identifier", resourceIdentifier.Identifier()),
	)
	_, err := s3Storage.client.DeleteObject(s3Ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s3Storage.configuration.BUCKET),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		span.RecordError(err)
	}
	tracer.SafeEndSpan(span)

	if err != nil {
		logger.Error(ctx, "failed to delete s3 object", "resource.identifier", resourceIdentifier.Identifier(), "s3.key", objectKey, "error", err)
		return err
	}

	logger.Debug(ctx, "successfully deleted s3 object", "resource.identifier", resourceIdentifier.Identifier(), "s3.key", objectKey)
	return nil
}

func (s3Storage *Storage) ResourceExists(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (bool, error) {
	objectKey := resourceIdentifier.ToUniqueFilename()
	logger.Debug(ctx, "checking if s3 object exists", "resource.identifier", resourceIdentifier.Identifier(), "s3.key", objectKey)

	s3Ctx, span := tracer.StartDebugSpan(ctx, "s3.head_object")
	tracer.SetDebugAttributes(s3Ctx, span,
		attribute.String("s3.bucket", s3Storage.configuration.BUCKET),
		attribute.String("s3.key", objectKey),
		attribute.String("resource.identifier", resourceIdentifier.Identifier()),
	)
	_, err := s3Storage.client.HeadObject(s3Ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s3Storage.configuration.BUCKET),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		span.RecordError(err)
	}
	tracer.SafeEndSpan(span)

	if err != nil {
		var notFound *types.NotFound
		if errors.As(err, &notFound) {
			logger.Debug(ctx, "s3 object not found", "resource.identifier", resourceIdentifier.Identifier(), "s3.key", objectKey)
			return false, nil
		}
		logger.Error(ctx, "failed to check s3 object existence", "resource.identifier", resourceIdentifier.Identifier(), "s3.key", objectKey, "error", err)
		return false, err
	}

	logger.Debug(ctx, "s3 object exists", "resource.identifier", resourceIdentifier.Identifier(), "s3.key", objectKey)
	return true, nil
}

func NewStorage(configuration *configuration.S3Configuration) storage.Storage {
	cfg, err := buildConfig(context.Background(), configuration)
	if err != nil {
		panic(err)
	}

	client := s3.NewFromConfig(cfg, buildS3Options(configuration)...)

	uploader := manager.NewUploader(client, func(u *manager.Uploader) {
		u.PartSize = configuration.PART_UPLOAD_SIZE
		u.Concurrency = configuration.UPLOAD_CONCURRENCY
	})

	downloader := manager.NewDownloader(client, func(d *manager.Downloader) {
		d.Concurrency = configuration.DOWNLOAD_CONCURRENCY
	})

	return &Storage{
		client:        client,
		uploader:      uploader,
		downloader:    downloader,
		configuration: *configuration,
	}
}

func buildS3Options(configuration *configuration.S3Configuration) []func(*s3.Options) {
	var options []func(*s3.Options)
	options = applyUsePathStyle(configuration, options)
	return options
}

func applyUsePathStyle(configuration *configuration.S3Configuration, options []func(*s3.Options)) []func(*s3.Options) {
	if configuration.BASE_ENDPOINT != "" {
		options = append(options, func(o *s3.Options) {
			o.UsePathStyle = true
		})
	}
	return options
}

func buildConfig(ctx context.Context, configuration *configuration.S3Configuration) (aws.Config, error) {
	var options []func(*config.LoadOptions) error
	options = applyRegion(configuration, options)
	options = applyEndpoint(configuration, options)
	options = applyCredentials(configuration, options)

	return config.LoadDefaultConfig(ctx, options...)
}

func applyEndpoint(cfg *configuration.S3Configuration, options []func(*config.LoadOptions) error) []func(*config.LoadOptions) error {
	if cfg.BASE_ENDPOINT != "" {
		options = append(options, config.WithBaseEndpoint(cfg.BASE_ENDPOINT))
	}
	return options
}

func applyCredentials(cfg *configuration.S3Configuration, options []func(*config.LoadOptions) error) []func(*config.LoadOptions) error {
	if cfg.ACCESS_KEY != "" && cfg.SECRET_KEY != "" {
		options = append(options, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.ACCESS_KEY, cfg.SECRET_KEY, ""),
		))
	}
	return options
}

func applyRegion(cfg *configuration.S3Configuration, options []func(*config.LoadOptions) error) []func(*config.LoadOptions) error {
	if cfg.REGION != "" {
		options = append(options, config.WithRegion(cfg.REGION))
	}
	return options
}

func (s3Storage *Storage) GetName() string {
	return "s3"
}

func (s3Storage *Storage) SaveResource(ctx context.Context, resource *resource.Resource) error {
	objectKey := resource.Identifier.ToUniqueFilename()
	reader := s3Storage.createResourceReader(ctx, resource)
	logger.Debug(ctx, "trying to upload s3 object", "resource.identifier", resource.Identifier.Identifier(), "s3.key", objectKey)

	s3Ctx, span := tracer.StartDebugSpan(ctx, "s3.put_object")
	tracer.SetDebugAttributes(s3Ctx, span,
		attribute.String("s3.bucket", s3Storage.configuration.BUCKET),
		attribute.String("s3.key", objectKey),
		attribute.String("resource.identifier", resource.Identifier.Identifier()),
		attribute.Int64("s3.part_size", s3Storage.configuration.PART_UPLOAD_SIZE),
		attribute.Int("s3.concurrency", s3Storage.configuration.UPLOAD_CONCURRENCY),
	)
	_, err := s3Storage.uploader.Upload(s3Ctx, &s3.PutObjectInput{
		Bucket: aws.String(s3Storage.configuration.BUCKET),
		Key:    aws.String(objectKey),
		Body:   reader,
	})
	if err != nil {
		span.RecordError(err)
	}
	tracer.SafeEndSpan(span)

	if err != nil {
		logger.Error(ctx, "failed to upload s3 object", "resource.identifier", resource.Identifier.Identifier(), "s3.key", objectKey, "error", err)
		return err
	}

	logger.Debug(ctx, "successfully uploaded s3 object", "resource.identifier", resource.Identifier.Identifier(), "s3.key", objectKey)
	return nil
}

func (s3Storage *Storage) createResourceReader(ctx context.Context, resource *resource.Resource) io.Reader {
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
