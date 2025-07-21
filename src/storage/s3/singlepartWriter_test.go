package s3

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestShouldCallPutObjectWithContextOnClose(t *testing.T) {
	ctx := context.Background()
	client := new(FakeClient)
	client.On("PutObjectContext", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
	singlepartWriter := NewSinglepartWriterWithContext(ctx, client, "test-bucket", "test-key", &bytes.Buffer{})

	singlepartWriter.Close()

	client.AssertCalled(t, "PutObjectContext", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}
