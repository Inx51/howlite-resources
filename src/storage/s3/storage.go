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
	"github.com/inx51/howlite/resources/config"
	"github.com/inx51/howlite/resources/resource"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type S3 struct {
	logger         *slog.Logger
	bucket         string
	prefix         string
	tracer         trace.Tracer
	client         AwsS3Client
	uploadStrategy string
	partUploadSize int
}

func NewStorage(config config.S3Configuration, logger *slog.Logger) (*S3, error) {
	tracer := otel.Tracer("S3")

	creds := aws.NewCredentialsCache(awscred.NewStaticCredentialsProvider(
		config.ACCESS_KEY,
		config.SECRET_KEY,
		"",
	))

	s3config := aws.Config{
		Region:      config.REGION,
		Credentials: creds,
	}

	if config.ENDPOINT != "" {
		s3config.BaseEndpoint = &config.ENDPOINT
	}

	s3client := NewAwsS3Client(s3config)

	return &S3{
		bucket:         config.BUCKET,
		prefix:         config.PREFIX,
		client:         *s3client,
		logger:         logger,
		tracer:         tracer,
		uploadStrategy: config.UPLOAD_STRATEGY,
		partUploadSize: config.MULTIPART_PART_UPLOAD_SIZE,
	}, nil
}

func (fileSystem *S3) GetName() string {
	return "S3"
}

func (s3 *S3) RemoveResourceContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) error {

	key := s3.resourceKey(resourceIdentifier)

	s3.logger.DebugContext(ctx, "Trying to remove resource file", "resourceIdentifier", resourceIdentifier.Value, "key", key)

	err := s3.client.DeleteObject(ctx,
		&s3.bucket,
		key)

	if err != nil {

		s3.logger.DebugContext(ctx, "Removing resource file failed with unhandled error", "resourceIdentifier", resourceIdentifier.Value, "key", key, "error", err)
		return err
	}

	s3.logger.InfoContext(ctx, "Successfully removed resource file", "resourceIdentifier", resourceIdentifier.Value, "key", key)
	return nil
}

func (s3 *S3) NewResourceWriterContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (io.WriteCloser, error) {
	key := *s3.resourceKey(resourceIdentifier)

	switch strings.ToLower(s3.uploadStrategy) {
	case "multipart":
		return NewMultipartWriter(
			&ctx,
			&s3.bucket,
			&key,
			&s3.client,
			s3.logger,
			s3.partUploadSize)
	case "singlepart":
		return &singlePartWriter{
			ctx:    &ctx,
			bucket: &s3.bucket,
			key:    &key,
			client: &s3.client,
			logger: s3.logger,
			buffer: &bytes.Buffer{},
		}, nil
	default:
		s3.logger.ErrorContext(ctx, "Invalid upload strategy", "strategy", s3.uploadStrategy)
		return nil, errors.New("invalid upload strategy: " + s3.uploadStrategy)
	}
}

func (s3 *S3) ResourceExistsContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (bool, error) {
	err := s3.client.HeadObject(ctx,
		&s3.bucket,
		s3.resourceKey(resourceIdentifier))

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

	output, err := s3.client.GetObject(ctx,
		&s3.bucket,
		key)

	if err != nil {
		s3.logger.ErrorContext(ctx, "Failed to get resource file", "resourceIdentifier", resourceIdentifier.Value, "key", key, "error", err)
		return nil, err
	}

	return output.Body, nil
}

func (s3 *S3) resourceKey(resourceIdentifier *resource.ResourceIdentifier) *string {
	// Use the same identifier logic as filesystem, but as S3 key
	id := identifierToStringVersion(resourceIdentifier)
	if s3.prefix != "" {
		id = strings.TrimRight(s3.prefix, "/") + "/" + id
	}
	id = id + ".bin"
	return &id
}

func identifierToStringVersion(resourceIdentifier *resource.ResourceIdentifier) string {
	var encBytes = md5.Sum([]byte(*resourceIdentifier.Value))
	return base64.URLEncoding.EncodeToString(encBytes[:])
}
