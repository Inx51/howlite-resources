package s3

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	awscred "github.com/aws/aws-sdk-go-v2/credentials"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/inx51/howlite/resources/config"
	"github.com/inx51/howlite/resources/resource"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type S3 struct {
	logger         *slog.Logger
	Bucket         string
	Prefix         string
	tracer         trace.Tracer
	client         *awss3.Client
	partUploadSize int
}

func NewStorage(config config.S3Configuration, logger *slog.Logger) *S3 {
	tracer := otel.Tracer("S3")

	creds := aws.NewCredentialsCache(awscred.NewStaticCredentialsProvider(
		config.ACCESS_KEY,
		config.SECRET_KEY,
		"",
	))

	endpoint := config.ENDPOINT
	s3config := aws.Config{
		Region:       config.REGION,
		Credentials:  creds,
		BaseEndpoint: &endpoint,
	}

	s3client := awss3.NewFromConfig(s3config)

	return &S3{
		Bucket:         config.BUCKET,
		Prefix:         config.PREFIX,
		client:         s3client,
		logger:         logger,
		tracer:         tracer,
		partUploadSize: config.PART_UPLOAD_SIZE,
	}
}

func (fileSystem *S3) GetName() string {
	return "S3"
}

func (s3 *S3) RemoveResourceContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) error {

	key := s3.resourceKey(resourceIdentifier)

	s3.logger.DebugContext(ctx, "Trying to remove resource file", "resourceIdentifier", resourceIdentifier.Value, "key", key)

	_, err := s3.client.DeleteObject(ctx, &awss3.DeleteObjectInput{
		Bucket: &s3.Bucket,
		Key:    key,
	})

	if err != nil {

		s3.logger.DebugContext(ctx, "Removing resource file failed with unhandled error", "resourceIdentifier", resourceIdentifier.Value, "key", key, "error", err)
		return err
	}

	s3.logger.InfoContext(ctx, "Successfully removed resource file", "resourceIdentifier", resourceIdentifier.Value, "key", key)
	return nil
}

func (s3 *S3) NewResourceWriterContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (io.WriteCloser, error) {
	key := *s3.resourceKey(resourceIdentifier)

	output, err := s3.client.CreateMultipartUpload(ctx, &awss3.CreateMultipartUploadInput{
		Bucket: &s3.Bucket,
		Key:    &key,
	})

	if err != nil {
		s3.logger.ErrorContext(ctx, "Failed to create multipart upload", "resourceIdentifier", resourceIdentifier.Value, "key", key, "error", err)
		return nil, err
	}

	return &multipartWriter{
		ctx:      &ctx,
		bucket:   &s3.Bucket,
		key:      &key,
		client:   s3.client,
		uploadId: output.UploadId,
		logger:   s3.logger,
		buffer:   &bytes.Buffer{},
		partSize: s3.partUploadSize,
	}, nil
}

func (s3 *S3) ResourceExistsContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (bool, error) {
	_, err := s3.client.HeadObject(ctx, &awss3.HeadObjectInput{
		Bucket: &s3.Bucket,
		Key:    s3.resourceKey(resourceIdentifier),
	})

	if err != nil {
		var responseError *awshttp.ResponseError
		if errors.As(err, &responseError) && responseError.ResponseError.HTTPStatusCode() == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s3 *S3) GetResourceContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (io.ReadCloser, error) {

	key := s3.resourceKey(resourceIdentifier)

	output, err := s3.client.GetObject(ctx, &awss3.GetObjectInput{
		Bucket: &s3.Bucket,
		Key:    key,
	})
	if err != nil {
		s3.logger.ErrorContext(ctx, "Failed to get resource file", "resourceIdentifier", resourceIdentifier.Value, "key", key, "error", err)
		return nil, err
	}

	return output.Body, nil
}

func (s3 *S3) resourceKey(resourceIdentifier *resource.ResourceIdentifier) *string {
	// Use the same identifier logic as filesystem, but as S3 key
	id := identifierToStringVersion(resourceIdentifier)
	if s3.Prefix != "" {
		id = strings.TrimRight(s3.Prefix, "/") + "/" + id
	}
	id = id + ".bin"
	return &id
}

func identifierToStringVersion(resourceIdentifier *resource.ResourceIdentifier) string {
	var encBytes = md5.Sum([]byte(*resourceIdentifier.Value))
	return base64.URLEncoding.EncodeToString(encBytes[:])
}
