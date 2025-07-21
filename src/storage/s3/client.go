package s3

import (
	"context"
	"io"
)

type Client interface {
	RemoveObjectContext(ctx context.Context, bucket *string, key *string) error
	GetObjectContext(ctx context.Context, bucket *string, key *string) (io.ReadCloser, error)
	ObjectExistsContext(ctx context.Context, bucket *string, key *string) (bool, error)
	PutObjectContext(ctx context.Context, bucket *string, key *string, body io.Reader) error
	UploadPartContext(ctx context.Context, bucket *string, key *string, partNumber int32, body io.Reader) (string, error)
	CreateMultipartUploadContext(ctx context.Context, bucket *string, key *string) (*string, error)
	AbortMultipartUploadContext(ctx context.Context, bucket *string, key *string, uploadId *string) error
	CompleteMultipartUploadContext(ctx context.Context, bucket *string, key *string, parts []CompletedPart) error
}
