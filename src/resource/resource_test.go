//go:build unit

package resource_test

import (
	"bytes"
	"encoding/binary"
	"io"
	"strings"
	"testing"

	"github.com/inx51/howlite-resources/resource"
	"github.com/vmihailenco/msgpack/v5"
)

func TestNewResourceShouldInitializeCorrectly(t *testing.T) {
	identifier := resource.NewResourceIdentifier("test-resource")
	body := io.NopCloser(strings.NewReader("test body content"))

	res := resource.NewResource(identifier, &body)

	if res.Body == nil {
		t.Fatal("Expected body to be set")
	}
	if len(*res.Headers.Headers()) != 0 {
		t.Fatal("Expected empty headers")
	}
}

func TestNewResourceShouldSetIdentifierCorrectly(t *testing.T) {
	identifier := resource.NewResourceIdentifier("test-resource")
	body := io.NopCloser(strings.NewReader("content"))

	res := resource.NewResource(identifier, &body)

	if res.Identifier.Identifier() != "test-resource" {
		t.Fatalf("Expected identifier 'test-resource', got %s", res.Identifier.Identifier())
	}
}

func TestWriteShouldWriteContentToWriter(t *testing.T) {
	identifier := resource.NewResourceIdentifier("test")
	bodyContent := "test content"
	body := io.NopCloser(strings.NewReader(bodyContent))
	res := resource.NewResource(identifier, &body)
	var buf bytes.Buffer
	writer := &testWriteCloser{&buf}

	err := res.Write(writer)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !strings.Contains(string(buf.Bytes()), bodyContent) {
		t.Fatal("Expected body content in written data")
	}
}

func TestLoadResourceShouldLoadWithHeaders(t *testing.T) {
	identifier := resource.NewResourceIdentifier("test")
	headers := map[string][]string{"Type": {"json"}}
	reader := createTestReader(headers, "body")

	res, err := resource.LoadResource(identifier, reader)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	loadedHeaders := *res.Headers.Headers()
	if len(loadedHeaders) != 1 || loadedHeaders["Type"][0] != "json" {
		t.Fatal("Expected headers to be loaded correctly")
	}
}

func TestLoadResourceShouldLoadWithoutHeaders(t *testing.T) {
	identifier := resource.NewResourceIdentifier("test")
	reader := createTestReader(nil, "body")

	res, err := resource.LoadResource(identifier, reader)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(*res.Headers.Headers()) != 0 {
		t.Fatal("Expected empty headers")
	}
}

func createTestReader(headers map[string][]string, body string) io.ReadCloser {
	var buf bytes.Buffer
	if headers != nil {
		packed, _ := msgpack.Marshal(headers)
		l := make([]byte, 8)
		binary.LittleEndian.PutUint64(l, uint64(len(packed)))
		buf.Write(l)
		buf.Write(packed)
	} else {
		l := make([]byte, 8)
		buf.Write(l)
	}
	buf.WriteString(body)
	return io.NopCloser(bytes.NewReader(buf.Bytes()))
}

type testWriteCloser struct {
	*bytes.Buffer
}

func (twc *testWriteCloser) Close() error {
	return nil
}
