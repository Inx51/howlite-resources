package s3

import (
	"bytes"
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShouldCallCreateMultipartUploadContextWhenWritingWritingFirstBytes(t *testing.T) {
	bucket := "Test"
	key := "Test"
	uploadId := "Test"
	etag := "Test"
	ctx := context.Background()
	client := new(FakeClient)
	multipartWriter := NewMultipartWriterWithContext(ctx, client, bucket, key, 1, &bytes.Buffer{})
	client.On("UploadPartContext", mock.Anything, mock.AnythingOfType("*string"), mock.AnythingOfType("*string"), mock.AnythingOfType("int32"), mock.AnythingOfType("*bytes.Reader")).Return(etag, nil)
	client.On("CreateMultipartUploadContext", mock.Anything, mock.AnythingOfType("*string"), mock.AnythingOfType("*string")).Return(&uploadId, nil)

	multipartWriter.Write(make([]byte, 1))

	client.AssertCalled(t, "CreateMultipartUploadContext", mock.Anything, &bucket, &key)
}

func TestShouldCallCompleteMultipartUploadContextWhenClosing(t *testing.T) {
	bucket := "Test"
	key := "Test"
	ctx := context.Background()
	client := new(FakeClient)
	multipartWriter := NewMultipartWriterWithContext(ctx, client, bucket, key, 1, &bytes.Buffer{})
	client.On("CompleteMultipartUploadContext", mock.Anything, mock.AnythingOfType("*string"), mock.AnythingOfType("*string"), mock.AnythingOfType("[]s3.CompletedPart")).Return(nil)

	multipartWriter.Close()

	client.AssertCalled(t, "CompleteMultipartUploadContext", mock.Anything, &bucket, &key, mock.AnythingOfType("[]s3.CompletedPart"))
}

func TestShouldResetBufferWhenBytesOfPartSizeWritten(t *testing.T) {
	bucket := "Test"
	key := "Test"
	uploadId := "Test"
	etag := "Test"
	bytesWritten := 17
	buffer := bytes.Buffer{}
	ctx := context.Background()
	client := new(FakeClient)
	client.On("CreateMultipartUploadContext", mock.Anything, mock.AnythingOfType("*string"), mock.AnythingOfType("*string")).Return(&uploadId, nil)
	client.On("UploadPartContext", mock.Anything, mock.AnythingOfType("*string"), mock.AnythingOfType("*string"), mock.AnythingOfType("int32"), mock.AnythingOfType("*bytes.Reader")).Return(etag, nil)
	multipartWriter := NewMultipartWriterWithContext(ctx, client, bucket, key, 15, &buffer)

	multipartWriter.Write(make([]byte, bytesWritten))

	assert.Equal(t, 0, buffer.Len())
}

func TestShouldCallUploadPartContextMultipleTimesAndNotMissCallForRemaningBytes(t *testing.T) {
	bucket := "Test"
	key := "Test"
	uploadId := "Test"
	etag := "Test"
	bytesWritten := 112
	bytesInPart := 10
	parts := math.Ceil(float64(bytesWritten) / float64(bytesInPart))
	buffer := bytes.Buffer{}
	ctx := context.Background()
	client := new(FakeClient)
	client.On("CreateMultipartUploadContext", mock.Anything, mock.AnythingOfType("*string"), mock.AnythingOfType("*string")).Return(&uploadId, nil)
	client.On("UploadPartContext", mock.Anything, mock.AnythingOfType("*string"), mock.AnythingOfType("*string"), mock.AnythingOfType("int32"), mock.AnythingOfType("*bytes.Reader")).Return(etag, nil)
	multipartWriter := NewMultipartWriterWithContext(ctx, client, bucket, key, bytesInPart, &buffer)

	multipartWriter.Write(make([]byte, bytesWritten))

	client.AssertNumberOfCalls(t, "UploadPartContext", int(math.Ceil(parts)))
	//ensure that the buffer is empty
	assert.Equal(t, 0, buffer.Len())
}

func TestShouldCallUploadPartContextMultipleTimesForEvenPartSizeAndBytes(t *testing.T) {
	bucket := "Test"
	key := "Test"
	uploadId := "Test"
	etag := "Test"
	bytesWritten := 100
	bytesInPart := 10
	buffer := bytes.Buffer{}
	ctx := context.Background()
	client := new(FakeClient)
	client.On("CreateMultipartUploadContext", mock.Anything, mock.AnythingOfType("*string"), mock.AnythingOfType("*string")).Return(&uploadId, nil)
	client.On("UploadPartContext", mock.Anything, mock.AnythingOfType("*string"), mock.AnythingOfType("*string"), mock.AnythingOfType("int32"), mock.AnythingOfType("*bytes.Reader")).Return(etag, nil)
	multipartWriter := NewMultipartWriterWithContext(ctx, client, bucket, key, bytesInPart, &buffer)

	multipartWriter.Write(make([]byte, bytesWritten))

	client.AssertNumberOfCalls(t, "UploadPartContext", 10)
	//ensure that the buffer is empty
	assert.Equal(t, 0, buffer.Len())
}

func TestShouldCallUploadPartContextOneTimeIfBytesAreFeverThanPartSize(t *testing.T) {
	bucket := "Test"
	key := "Test"
	uploadId := "Test"
	etag := "Test"
	bytesWritten := 20
	bytesInPart := 100
	buffer := bytes.Buffer{}
	ctx := context.Background()
	client := new(FakeClient)
	client.On("CreateMultipartUploadContext", mock.Anything, mock.AnythingOfType("*string"), mock.AnythingOfType("*string")).Return(&uploadId, nil)
	client.On("UploadPartContext", mock.Anything, mock.AnythingOfType("*string"), mock.AnythingOfType("*string"), mock.AnythingOfType("int32"), mock.AnythingOfType("*bytes.Reader")).Return(etag, nil)
	multipartWriter := NewMultipartWriterWithContext(ctx, client, bucket, key, bytesInPart, &buffer)

	multipartWriter.Write(make([]byte, bytesWritten))

	client.AssertNumberOfCalls(t, "UploadPartContext", 1)
	//ensure that the buffer is empty
	assert.Equal(t, 0, buffer.Len())
}
