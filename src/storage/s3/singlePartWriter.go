package s3

import (
	"bytes"
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type singlePartWriter struct {
	ctx    *context.Context
	bucket *string
	key    *string
	client *s3.Client
	logger *slog.Logger
	buffer *bytes.Buffer
}

func (writer *singlePartWriter) Write(p []byte) (int, error) {
	return writer.buffer.Write(p)
}

func (writer *singlePartWriter) Close() error {
	_, err := writer.client.PutObject(*writer.ctx, &s3.PutObjectInput{
		Bucket: writer.bucket,
		Key:    writer.key,
		Body:   bytes.NewReader(writer.buffer.Bytes()),
	})

	if err != nil {
		writer.logger.ErrorContext(*writer.ctx, "Failed to upload single part", "error", err)
		return err
	}

	return nil
}
