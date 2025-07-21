package s3

import (
	"bytes"
	"context"
)

type multipartWriter struct {
	bucketName string
	key        string
	client     Client
	ctx        context.Context
	parts      []CompletedPart
	uploadId   string
	partSize   int
	buffer     *bytes.Buffer
}

func NewMultipartWriterWithContext(ctx context.Context, client Client, bucketName string, key string, partSize int, buffer *bytes.Buffer) *multipartWriter {
	return &multipartWriter{
		bucketName: bucketName,
		key:        key,
		client:     client,
		ctx:        ctx,
		parts:      []CompletedPart{},
		uploadId:   "",
		partSize:   partSize,
		buffer:     buffer,
	}
}

// TODO: if failed to write parts or complete, should we abort the upload?
func (mw *multipartWriter) Write(p []byte) (n int, err error) {
	if mw.uploadId == "" {
		uploadId, err := mw.client.CreateMultipartUploadContext(mw.ctx, &mw.bucketName, &mw.key)
		if err != nil {
			return 0, err
		}
		mw.uploadId = *uploadId
	}

	written := 0
	for written < len(p) {
		spaceLeft := mw.partSize - mw.buffer.Len()
		toWrite := len(p) - written
		if toWrite > spaceLeft {
			toWrite = spaceLeft
		}
		mw.buffer.Write(p[written : written+toWrite])
		written += toWrite

		if mw.buffer.Len() == mw.partSize {
			if err := mw.uploadPart(); err != nil {
				return written, err
			}
		}
	}

	if mw.buffer.Len() > 0 {
		if err := mw.uploadPart(); err != nil {
			return written, err
		}
	}

	return written, nil
}

func (multipartWriter *multipartWriter) uploadPart() error {
	if multipartWriter.buffer.Len() == 0 {
		return nil
	}

	partNumber := int32(len(multipartWriter.parts))
	partData := bytes.NewReader(multipartWriter.buffer.Bytes())
	etag, err := multipartWriter.client.UploadPartContext(multipartWriter.ctx, &multipartWriter.bucketName, &multipartWriter.key, partNumber, partData)
	if err != nil {
		return err
	}

	multipartWriter.parts = append(multipartWriter.parts, CompletedPart{
		ETag:       &etag,
		PartNumber: &partNumber,
	})
	multipartWriter.buffer.Reset()

	return nil
}

func (multipartWriter *multipartWriter) Close() error {
	return multipartWriter.client.CompleteMultipartUploadContext(multipartWriter.ctx, &multipartWriter.bucketName, &multipartWriter.key, multipartWriter.parts)
}
