package s3

import (
	"bytes"
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/inx51/howlite/resources/errors"
)

type multipartWriter struct {
	ctx      *context.Context
	bucket   *string
	key      *string
	client   *s3.Client
	uploadId *string
	logger   *slog.Logger
	parts    []types.CompletedPart
	buffer   *bytes.Buffer
	partSize int
}

func (writer *multipartWriter) Write(p []byte) (int, error) {

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

func (writer *multipartWriter) Close() error {

	writer.completeMultipartUpload()

	return nil
}

func (writer *multipartWriter) uploadPart() error {

	partNumber := int32(len(writer.parts) + 1)

	output, err := writer.client.UploadPart(*writer.ctx, &s3.UploadPartInput{
		Bucket:     writer.bucket,
		Key:        writer.key,
		UploadId:   writer.uploadId,
		Body:       bytes.NewReader(writer.buffer.Bytes()),
		PartNumber: &partNumber,
	})

	writer.buffer = bytes.NewBuffer(nil)

	if err != nil {
		writer.logger.ErrorContext(*writer.ctx, "Failed to upload part", "error", err, "partNumber", partNumber)
		abortErr := writer.abortMultipartUpload()

		if abortErr != nil {
			return errors.NewAggregatedError([]error{err, abortErr})
		}

		return err
	}

	writer.parts = append(writer.parts, types.CompletedPart{
		ETag:       output.ETag,
		PartNumber: &partNumber,
	})

	return nil
}

func (writer *multipartWriter) completeMultipartUpload() error {

	if writer.buffer.Len() > 0 {
		err := writer.uploadPart()
		if err != nil {
			return err
		}
	}

	_, err := writer.client.CompleteMultipartUpload(
		*writer.ctx,
		&s3.CompleteMultipartUploadInput{
			Bucket:   writer.bucket,
			Key:      writer.key,
			UploadId: writer.uploadId,
			MultipartUpload: &types.CompletedMultipartUpload{
				Parts: writer.parts},
		},
	)

	if err != nil {
		writer.logger.ErrorContext(*writer.ctx, "Failed to complete multipart upload", "error", err)
		abortErr := writer.abortMultipartUpload()

		if abortErr != nil {
			return errors.NewAggregatedError([]error{err, abortErr})
		}

		return err
	}

	return nil
}

func (writer *multipartWriter) abortMultipartUpload() error {

	_, err := writer.client.AbortMultipartUpload(*writer.ctx, &s3.AbortMultipartUploadInput{
		Bucket:   writer.bucket,
		Key:      writer.key,
		UploadId: writer.uploadId,
	})

	if err != nil {
		writer.logger.ErrorContext(*writer.ctx, "Failed to abort multipart upload", "error", err)
	}

	return nil
}
