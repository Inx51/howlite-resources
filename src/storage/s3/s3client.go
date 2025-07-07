package s3

import (
	"bytes"
	"context"
)

type S3Client interface {
	CreateMultipartUpload(ctx context.Context, bucket *string, key *string) (uploadId string, err error)
	UploadPart(ctx context.Context, bucket *string, key *string, uploadId *string, body *bytes.Reader, partNumber *int32) (etag *string, err error)
	CompleteMultipartUpload(ctx context.Context, bucket *string, key *string, uploadId *string, parts *[]CompletedPart) error
	AbortMultipartUpload(ctx context.Context, bucket *string, key *string, uploadId string) error
	PutObject(ctx context.Context, bucket *string, key *string, body bytes.Reader) error
	HeadObject(ctx, bucket *string, key *string) (err error)
	GetObject(ctx context.Context, bucket *string, key *string) (body *bytes.Reader, err error)
}
