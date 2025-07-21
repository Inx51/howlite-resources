package s3

import (
	"context"
	"io"

	"github.com/stretchr/testify/mock"
)

type FakeClient struct {
	mock.Mock
}

func (client *FakeClient) RemoveObjectContext(ctx context.Context, bucket *string, key *string) error {
	args := client.Called(ctx, bucket, key)
	return args.Error(0)
}

func (client *FakeClient) GetObjectContext(ctx context.Context, bucket *string, key *string) (io.ReadCloser, error) {
	args := client.Called(ctx, bucket, key)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (client *FakeClient) ObjectExistsContext(ctx context.Context, bucket *string, key *string) (bool, error) {
	args := client.Called(ctx, bucket, key)
	return args.Bool(0), args.Error(1)
}

func (client *FakeClient) PutObjectContext(ctx context.Context, bucket *string, key *string, body io.Reader) error {
	args := client.Called(ctx, bucket, key, body)
	return args.Error(0)
}

func (client *FakeClient) UploadPartContext(ctx context.Context, bucket *string, key *string, partNumber int32, body io.Reader) (string, error) {
	args := client.Called(ctx, bucket, key, partNumber, body)
	return args.String(0), args.Error(1)
}

func (client *FakeClient) CreateMultipartUploadContext(ctx context.Context, bucket *string, key *string) (*string, error) {
	args := client.Called(ctx, bucket, key)
	return args.Get(0).(*string), args.Error(1)
}

func (client *FakeClient) AbortMultipartUploadContext(ctx context.Context, bucket *string, key *string, uploadId *string) error {
	args := client.Called(ctx, bucket, key, uploadId)
	return args.Error(0)
}

func (client *FakeClient) CompleteMultipartUploadContext(ctx context.Context, bucket *string, key *string, parts []CompletedPart) error {
	args := client.Called(ctx, bucket, key, parts)
	return args.Error(0)
}
