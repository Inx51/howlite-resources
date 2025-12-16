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
	"github.com/inx51/howlite-resources/resource"
	"github.com/inx51/howlite-resources/storage"
)

type Storage struct {
	client        *s3.Client
	uploader      *manager.Uploader
	downloader    *manager.Downloader
	configuration configuration.S3Configuration
}

func (s3Storage *Storage) GetResource(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (*resource.Resource, error) {
	result, err := s3Storage.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s3Storage.configuration.BUCKET),
		Key:    aws.String(resourceIdentifier.ToUniqueFilename()),
	})
	if err != nil {
		return nil, err
	}

	resource, err := resource.LoadResource(resourceIdentifier, result.Body)

	return resource, err
}

func (s3Storage *Storage) RemoveResource(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) error {
	_, err := s3Storage.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s3Storage.configuration.BUCKET),
		Key:    aws.String(resourceIdentifier.ToUniqueFilename()),
	})

	return err
}

func (s3Storage *Storage) ResourceExists(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (bool, error) {
	_, err := s3Storage.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s3Storage.configuration.BUCKET),
		Key:    aws.String(resourceIdentifier.ToUniqueFilename()),
	})
	if err != nil {
		var notFound *types.NotFound
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func NewStorage(configuration *configuration.S3Configuration) storage.Storage {
	cfg, err := loadAWSConfig(context.Background(), configuration)
	if err != nil {
		panic(err)
	}

	client := s3.NewFromConfig(cfg)

	uploader := manager.NewUploader(client, func(u *manager.Uploader) {
		u.PartSize = int64(configuration.MULTIPART_PART_UPLOAD_SIZE)
		u.Concurrency = 5
	})

	downloader := manager.NewDownloader(client, func(d *manager.Downloader) {
		d.Concurrency = 5
	})

	return &Storage{
		client:        client,
		uploader:      uploader,
		downloader:    downloader,
		configuration: *configuration,
	}
}

func loadAWSConfig(ctx context.Context, cfg *configuration.S3Configuration) (aws.Config, error) {
	var options []func(*config.LoadOptions) error

	//TODO: Make sure to extract this into functions..
	if cfg.REGION != "" {
		options = append(options, config.WithRegion(cfg.REGION))
	}

	if cfg.ENDPOINT != "" {
		options = append(options, config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:               cfg.ENDPOINT,
					HostnameImmutable: true,
				}, nil
			}),
		))
	}

	if cfg.ACCESS_KEY != "" && cfg.SECRET_KEY != "" {
		options = append(options, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.ACCESS_KEY, cfg.SECRET_KEY, ""),
		))
	}

	return config.LoadDefaultConfig(ctx, options...)
}

func (s3Storage *Storage) GetName() string {
	return "s3"
}

func (s3Storage *Storage) SaveResource(ctx context.Context, resource *resource.Resource) error {
	reader := s3Storage.createResourceReader(resource)

	// Use manager.Uploader for automatic multipart uploads with parallelization
	_, err := s3Storage.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s3Storage.configuration.BUCKET),
		Key:    aws.String(resource.Identifier.ToUniqueFilename()),
		Body:   reader,
	})

	return err
}

func (s3Storage *Storage) createResourceReader(resource *resource.Resource) io.Reader {
	pipeReader, pipeWriter := io.Pipe()
	go func() {
		defer pipeWriter.Close()
		err := resource.Write(pipeWriter)
		if err != nil {
			pipeWriter.CloseWithError(err)
		}
	}()
	return pipeReader
}
