package s3

import (
	"bytes"
	"context"
	"errors"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	customerrors "github.com/inx51/howlite/resources/errors"
)

type MultipartWriter struct {
	ctx      *context.Context
	bucket   *string
	key      *string
	client   S3Client
	uploadId *string
	logger   *slog.Logger
	parts    []types.CompletedPart
	buffer   *bytes.Buffer
	partSize int
}

func NewMultipartWriter(
	ctx *context.Context,
	bucket *string,
	key *string,
	client S3Client,
	logger *slog.Logger,
	partSize int) (*MultipartWriter, error) {
	if partSize <= 0 {
		return nil, errors.New("MULTIPART_PART_UPLOAD_SIZE must be greater than 0")
	}

	uploadId, err := client.CreateMultipartUpload(*ctx, bucket, key)

	if err != nil {
		logger.ErrorContext(*ctx, "Failed to create multipart upload", "resourceIdentifier", key, "error", err)
		return nil, err
	}

	return &MultipartWriter{
		ctx:      ctx,
		bucket:   bucket,
		key:      key,
		client:   client,
		logger:   logger,
		uploadId: &uploadId,
		buffer:   &bytes.Buffer{},
		partSize: partSize,
	}, nil
}

func (writer *MultipartWriter) Write(p []byte) (int, error) {
	take := 0
	if writer.buffer.Len()+len(p) > writer.partSize {
		take = writer.partSize - writer.buffer.Len()
		writer.buffer.Write(p[:take])

		if len(p)-take > 0 {
			writer.buffer.Write(p[take:])
		}

		writer.uploadPart()

	} else {
		writer.buffer.Write(p)
	}

	return len(p), nil
}

func (writer *MultipartWriter) Close() error {

	writer.completeMultipartUpload()

	return nil
}

func (writer *MultipartWriter) uploadPart() error {

	partNumber := int32(len(writer.parts) + 1)

	etag, err := writer.client.UploadPart(*writer.ctx,
		writer.bucket,
		writer.key,
		writer.uploadId,
		bytes.NewReader(writer.buffer.Bytes()),
		&partNumber)

	writer.buffer = bytes.NewBuffer(nil)

	if err != nil {
		writer.logger.ErrorContext(*writer.ctx, "Failed to upload part", "error", err, "partNumber", partNumber)
		abortErr := writer.abortMultipartUpload()

		if abortErr != nil {
			return customerrors.NewAggregatedError([]error{err, abortErr})
		}

		return err
	}

	writer.parts = append(writer.parts, types.CompletedPart{
		ETag:       etag,
		PartNumber: &partNumber,
	})

	return nil
}

func (writer *MultipartWriter) completeMultipartUpload() error {

	if writer.buffer.Len() > 0 {
		err := writer.uploadPart()
		if err != nil {
			return err
		}
	}

	completedParts := make([]CompletedPart, len(writer.parts))
	for i, part := range writer.parts {
		completedParts[i] = NewCompletedPart(part.ETag, part.PartNumber)
	}

	err := writer.client.CompleteMultipartUpload(*writer.ctx,
		writer.bucket,
		writer.key,
		writer.uploadId,
		&completedParts)

	if err != nil {
		writer.logger.ErrorContext(*writer.ctx, "Failed to complete multipart upload", "error", err)
		abortErr := writer.abortMultipartUpload()

		if abortErr != nil {
			return customerrors.NewAggregatedError([]error{err, abortErr})
		}

		return err
	}

	return nil
}

func (writer *MultipartWriter) abortMultipartUpload() error {

	err := writer.client.AbortMultipartUpload(*writer.ctx,
		writer.bucket,
		writer.key,
		*writer.uploadId)

	if err != nil {
		writer.logger.ErrorContext(*writer.ctx, "Failed to abort multipart upload", "error", err)
	}

	return nil
}
