package s3

import (
	"bytes"
	"context"
	"log/slog"
)

type singlePartWriter struct {
	ctx    *context.Context
	bucket *string
	key    *string
	client S3Client
	logger *slog.Logger
	buffer *bytes.Buffer
}

func (writer *singlePartWriter) Write(p []byte) (int, error) {
	return writer.buffer.Write(p)
}

func (writer *singlePartWriter) Close() error {
	err := writer.client.PutObject(*writer.ctx,
		writer.bucket,
		writer.key,
		*bytes.NewReader(writer.buffer.Bytes()))

	if err != nil {
		writer.logger.ErrorContext(*writer.ctx, "Failed to upload single part", "error", err)
		return err
	}

	return nil
}
