package s3

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type AwsS3Client struct {
	client *s3.Client
}

func NewAwsS3Client(cfg aws.Config) *AwsS3Client {
	s3Client := s3.NewFromConfig(cfg)
	return &AwsS3Client{
		client: s3Client,
	}
}

func (a *AwsS3Client) CreateMultipartUpload(ctx context.Context, bucket *string, key *string) (uploadId string, err error) {
	input := &s3.CreateMultipartUploadInput{
		Bucket: bucket,
		Key:    key,
	}
	output, err := a.client.CreateMultipartUpload(ctx, input)
	if err != nil {
		return "", err
	}
	if output.UploadId == nil {
		return "", nil
	}
	return *output.UploadId, nil
}

func (a *AwsS3Client) UploadPart(ctx context.Context, bucket *string, key *string, uploadId *string, body *bytes.Reader, partNumber *int32) (etag *string, err error) {
	input := &s3.UploadPartInput{
		Bucket:     bucket,
		Key:        key,
		UploadId:   uploadId,
		PartNumber: partNumber,
		Body:       body,
	}
	output, err := a.client.UploadPart(ctx, input)
	if err != nil {
		return nil, err
	}
	if output.ETag == nil {
		return nil, nil
	}
	return output.ETag, nil
}

func (a *AwsS3Client) CompleteMultipartUpload(ctx context.Context, bucket *string, key *string, uploadId *string, parts *[]CompletedPart) error {
	derefParts := *parts
	completedParts := make([]types.CompletedPart, len(derefParts))
	for i, part := range derefParts {
		completedParts[i] = types.CompletedPart{
			PartNumber: part.PartNumber,
			ETag:       part.ETag,
		}
	}
	input := &s3.CompleteMultipartUploadInput{
		Bucket:   bucket,
		Key:      key,
		UploadId: uploadId,
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: completedParts,
		},
	}
	_, err := a.client.CompleteMultipartUpload(ctx, input)
	return err
}

func (a *AwsS3Client) AbortMultipartUpload(ctx context.Context, bucket *string, key *string, uploadId string) error {
	input := &s3.AbortMultipartUploadInput{
		Bucket:   bucket,
		Key:      key,
		UploadId: &uploadId,
	}
	_, err := a.client.AbortMultipartUpload(ctx, input)
	return err
}

func (a *AwsS3Client) PutObject(ctx context.Context, bucket *string, key *string, body bytes.Reader) error {
	input := &s3.PutObjectInput{
		Bucket: bucket,
		Key:    key,
		Body:   &body,
	}
	_, err := a.client.PutObject(ctx, input)
	return err
}

func (a *AwsS3Client) HeadObject(ctx context.Context, bucket *string, key *string) error {
	input := &s3.HeadObjectInput{
		Bucket: bucket,
		Key:    key,
	}
	_, err := a.client.HeadObject(ctx, input)
	return err
}

func (a *AwsS3Client) GetObject(ctx context.Context, bucket *string, key *string) (*bytes.Reader, error) {
	input := &s3.GetObjectInput{
		Bucket: bucket,
		Key:    key,
	}
	output, err := a.client.GetObject(ctx, input)
	if err != nil {
		return nil, err
	}
	if output.Body == nil {
		return nil, nil
	}
	return output.Body, nil
}
