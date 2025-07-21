package s3

import (
	"bytes"
	"context"
)

type singlepartWriter struct {
	bucketName string
	key        string
	client     Client
	ctx        context.Context
	buffer     bytes.Buffer
}

func NewSinglepartWriterWithContext(ctx context.Context, client Client, bucketName string, key string, buffer *bytes.Buffer) *singlepartWriter {
	return &singlepartWriter{
		bucketName: bucketName,
		key:        key,
		client:     client,
		ctx:        ctx,
		buffer:     *buffer,
	}
}

func (singlepartWriter *singlepartWriter) Write(p []byte) (int, error) {
	return singlepartWriter.buffer.Write(p)
}

func (singlepartWriter *singlepartWriter) Close() error {
	data := bytes.NewReader(singlepartWriter.buffer.Bytes())
	return singlepartWriter.client.PutObjectContext(singlepartWriter.ctx, &singlepartWriter.bucketName, &singlepartWriter.key, data)
}
