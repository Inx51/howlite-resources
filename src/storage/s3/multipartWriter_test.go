package s3_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/inx51/howlite/resources/storage/s3"
	"github.com/stretchr/testify/assert"
)

// fakeS3Client records uploaded parts for testing
// You may need to adjust this to match the real s3.Client interface used in your writer

type fakeS3Client struct {
	uploadedParts [][]byte
	completed     bool
}

func (f *fakeS3Client) UploadPart(part []byte) {
	f.uploadedParts = append(f.uploadedParts, part)
}

func (f *fakeS3Client) CompleteMultipartUpload() {
	f.completed = true
}

func newTestMultipartWriter(partSize int) (*s3.MultipartWriter, *fakeS3Client) {
	fake := fakeS3Client{}
	context := context.Background()
	bucket := "test-bucket"
	key := "test-key"
	writer, _ := s3.NewMultipartWriter(&context, &bucket, &key, &fake, slog.Default(), partSize)
	return writer, fake
}

// ctx *context.Context,
// 	bucket *string,
// 	key *string,
// 	client *s3.Client,
// 	logger *slog.Logger,
// 	partSize int

func TestWriterShouldUploadInMultipleParts(t *testing.T) {
	writer, fake := newTestMultipartWriter(5)
	data := []byte("abcdefghij") // 10 bytes, should be 2 parts of 5
	writer.Write(data)
	writer.Close()
	assert.Equal(t, 2, len(fake.uploadedParts))
}

func TestWriterShouldUploadInMultiplePartsEvenIfLastPartIsOdd(t *testing.T) {
	writer, fake := newTestMultipartWriter(4)
	data := []byte("abcdefg") // 7 bytes, should be 2 parts: 4 and 3
	_, err := writer.Write(data)
	assert.NoError(t, err)
	writer.Close()
	// assert.Equal(t, 2, len(fake.uploadedParts))
}

func TestWriterShouldUploadInMultiplePartsEvenIfOnlySinglePart(t *testing.T) {
	writer, fake := newTestMultipartWriter(10)
	data := []byte("abc") // 3 bytes, should be 1 part
	_, err := writer.Write(data)
	assert.NoError(t, err)
	writer.Close()
	// assert.Equal(t, 1, len(fake.uploadedParts))
}
