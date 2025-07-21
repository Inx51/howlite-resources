package s3

import (
	"context"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/inx51/howlite/resources/config"
)

type AwsClient struct {
	client *awss3.Client
}

func NewAwsClient(config *config.S3Configuration) *AwsClient {

	creds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		config.ACCESS_KEY,
		config.SECRET_KEY,
		"",
	))

	s3config := aws.Config{
		Region:      config.REGION,
		Credentials: creds,
	}

	if config.ENDPOINT != "" {
		s3config.BaseEndpoint = &config.ENDPOINT
	}

	s3client := awss3.NewFromConfig(s3config)

	return &AwsClient{
		client: s3client,
	}
}

func (awsClient *AwsClient) RemoveObjectContext(ctx context.Context, bucket *string, key *string) error {
	_, err := awsClient.client.DeleteObject(ctx, &awss3.DeleteObjectInput{
		Bucket: bucket,
		Key:    key,
	})
	return err
}

func (awsClient *AwsClient) GetObjectContext(ctx context.Context, bucket *string, key *string) (io.ReadCloser, error) {
	output, err := awsClient.client.GetObject(ctx, &awss3.GetObjectInput{
		Bucket: bucket,
		Key:    key,
	})
	if err != nil {
		return nil, err
	}
	return output.Body, nil
}

func (awsClient *AwsClient) ObjectExistsContext(ctx context.Context, bucket *string, key *string) (bool, error) {
	_, err := awsClient.client.HeadObject(ctx, &awss3.HeadObjectInput{
		Bucket: bucket,
		Key:    key,
	})
	if err != nil {
		var notFoundErr *types.NotFound
		if ok := errors.As(err, &notFoundErr); ok {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (awsClient *AwsClient) PutObjectContext(ctx context.Context, bucket *string, key *string, body io.Reader) error {
	var reader io.Reader
	if body != nil {
		reader = body
	}
	_, err := awsClient.client.PutObject(ctx, &awss3.PutObjectInput{
		Bucket: bucket,
		Key:    key,
		Body:   reader,
	})
	return err
}

func (awsClient *AwsClient) UploadPartContext(ctx context.Context, bucket *string, key *string, partNumber int32, body io.Reader) (string, error) {
	output, err := awsClient.client.UploadPart(ctx, &awss3.UploadPartInput{
		Bucket:     bucket,
		Key:        key,
		PartNumber: &partNumber,
		Body:       body,
	})
	if err != nil {
		return "", err
	}
	return *output.ETag, nil
}

func (awsClient *AwsClient) CreateMultipartUploadContext(ctx context.Context, bucket *string, key *string) (*string, error) {
	output, err := awsClient.client.CreateMultipartUpload(ctx, &awss3.CreateMultipartUploadInput{
		Bucket: bucket,
		Key:    key,
	})
	return output.UploadId, err
}

func (awsClient *AwsClient) AbortMultipartUploadContext(ctx context.Context, bucket *string, key *string, uploadId *string) error {
	_, err := awsClient.client.AbortMultipartUpload(ctx, &awss3.AbortMultipartUploadInput{
		Bucket:   bucket,
		Key:      key,
		UploadId: uploadId,
	})
	return err
}

func (awsClient *AwsClient) CompleteMultipartUploadContext(ctx context.Context, bucket *string, key *string, parts []CompletedPart) error {
	completedParts := make([]types.CompletedPart, len(parts))
	for i, part := range parts {
		completedParts[i] = types.CompletedPart{
			ETag:       part.ETag,
			PartNumber: part.PartNumber,
		}
	}
	_, err := awsClient.client.CompleteMultipartUpload(ctx, &awss3.CompleteMultipartUploadInput{
		Bucket: bucket,
		Key:    key,
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: completedParts,
		},
	})
	return err
}
