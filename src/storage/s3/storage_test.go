package s3

import (
	"context"
	"log/slog"
	"testing"

	"github.com/inx51/howlite/resources/config"
	"github.com/inx51/howlite/resources/resource"
	"github.com/inx51/howlite/resources/testing/utilities/logging"
	"github.com/stretchr/testify/assert"
)

func TestNewResourceWriterContextShouldUseMultipartStrategy(t *testing.T) {
	identifier := "test-resource"
	ctx := context.Background()
	resourceIdentifier := &resource.ResourceIdentifier{Value: &identifier}
	config := config.S3Configuration{
		UPLOAD_STRATEGY: "singlepart",
	}
	client := FakeClient{}
	logger := slog.New(slog.NewTextHandler(&logging.TestingLogWriter{}, nil))
	storage, _ := NewStorage(config, &client, logger)

	writer, _ := storage.NewResourceWriterContext(ctx, resourceIdentifier)

	_, ok := writer.(*singlepartWriter)

	assert.True(t, ok)
}

func TestNewResourceWriterContextShouldUseSinglepartStrategy(t *testing.T) {
	identifier := "test-resource"
	ctx := context.Background()
	resourceIdentifier := &resource.ResourceIdentifier{Value: &identifier}
	config := config.S3Configuration{
		UPLOAD_STRATEGY: "multipart",
	}
	client := FakeClient{}
	logger := slog.New(slog.NewTextHandler(&logging.TestingLogWriter{}, nil))
	storage, _ := NewStorage(config, &client, logger)

	writer, _ := storage.NewResourceWriterContext(ctx, resourceIdentifier)

	_, ok := writer.(*multipartWriter)

	assert.True(t, ok)
}

func TestNewResourceWriterContextShouldThrowUndefinedStrategyErrorIfUnimplementedStrategyIsConfigured(t *testing.T) {
	identifier := "test-resource"
	ctx := context.Background()
	resourceIdentifier := &resource.ResourceIdentifier{Value: &identifier}
	config := config.S3Configuration{
		UPLOAD_STRATEGY: "NOT_IMPLEMENTED",
	}
	client := FakeClient{}
	logger := slog.New(slog.NewTextHandler(&logging.TestingLogWriter{}, nil))
	storage, _ := NewStorage(config, &client, logger)

	_, err := storage.NewResourceWriterContext(ctx, resourceIdentifier)

	assert.Error(t, err)
}
