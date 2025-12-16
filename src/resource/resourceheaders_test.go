//go:build unit

package resource_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"strings"
	"testing"

	"github.com/inx51/howlite-resources/resource"
	"github.com/vmihailenco/msgpack/v5"
)

func TestAddShouldAddHeaders(t *testing.T) {
	ctx := context.Background()
	resourceHeaders := resource.NewResourceHeaders()
	resourceHeaders.Add(ctx, "A", []string{"1", "2"})
	resourceHeaders.Add(ctx, "B", []string{"3"})

	headers := *resourceHeaders.Headers()

	if len(headers) != 2 {
		t.Fatal("Headers() length mismatch, got", len(headers), "want 2")
	}
}

func TestAddShouldHaveValues(t *testing.T) {
	ctx := context.Background()
	resourceHeaders := resource.NewResourceHeaders()
	resourceHeaders.Add(ctx, "A", []string{"1", "2"})
	resourceHeaders.Add(ctx, "B", []string{"3"})

	headers := *resourceHeaders.Headers()

	if headers["A"][0] != "1" || headers["A"][1] != "2" {
		t.Fatalf("Headers A values mismatch, got %v, want [1 2]", headers["A"])
	}
	if headers["B"][0] != "3" {
		t.Fatalf("Headers B values mismatch, got %v, want [3]", headers["B"])
	}
}

func TestLoadHeadersShouldLoadHeadersFromReaderCloser(t *testing.T) {
	resourceHeaders := resource.NewResourceHeaders()

	originalHeaders := map[string][]string{"a": {"b", "c"}}
	packed, _ := msgpack.Marshal(originalHeaders)
	var buf bytes.Buffer
	l := make([]byte, 8)
	binary.LittleEndian.PutUint64(l, uint64(len(packed)))
	buf.Write(l)
	buf.Write(packed)
	rc := io.NopCloser(io.NopCloser(bytes.NewReader(buf.Bytes())))
	_ = resourceHeaders.LoadHeaders(rc)

	headers := *resourceHeaders.Headers()

	if len(headers) != len(originalHeaders) {
		t.Fatalf("LoadHeaders did not load correctly: got %v, want %v", headers, originalHeaders)
	}

	for k, v := range originalHeaders {
		hv, ok := headers[k]
		if !ok {
			t.Fatalf("LoadHeaders missing key: %s", k)
		}
		if strings.Join(v, ",") != strings.Join(hv, ",") {
			t.Fatalf("LoadHeaders value mismatch for key %s: got %v, want %v", k, hv, v)
		}
	}
}

func TestLoadHeadersShouldAsEmptyFromEmptyReaderCloser(t *testing.T) {
	resourceHeaders := resource.NewResourceHeaders()
	var buf bytes.Buffer
	l := make([]byte, 8)
	buf.Write(l)
	rc := io.NopCloser(io.NopCloser(bytes.NewReader(buf.Bytes())))
	_ = resourceHeaders.LoadHeaders(rc)

	headers := *resourceHeaders.Headers()

	if len(headers) != 0 {
		t.Fatalf("LoadHeaders did not load correctly: got %v, want empty", headers)
	}
}

func TestAddShouldBlockReservedHeaders(t *testing.T) {
	testCases := []struct {
		name         string
		headerName   string
		headerValues []string
		shouldBlock  bool
	}{
		{"content-length", "content-length", []string{"100"}, true},
		{"transfer-encoding", "transfer-encoding", []string{"chunked"}, true},
		{"connection", "connection", []string{"close"}, true},
		{"upgrade", "upgrade", []string{"websocket"}, true},
		{"server", "server", []string{"nginx/1.18"}, true},
		{"date", "date", []string{"Mon, 09 Dec 2025 12:00:00 GMT"}, true},
		{"trailer", "trailer", []string{"Expires"}, true},
		{"set-cookie", "set-cookie", []string{"session=abc123"}, true},
		{"set-cookie2", "set-cookie2", []string{"session=abc123"}, true},
		{"location", "location", []string{"https://example.com"}, true},
		{"retry-after", "retry-after", []string{"120"}, true},
		{"vary", "vary", []string{"Accept-Encoding"}, true},
		{"warning", "warning", []string{"110 anderson/1.3.37 \"Response is stale\""}, true},
		{"www-authenticate", "www-authenticate", []string{"Basic realm=\"Access\""}, true},
		{"proxy-authenticate", "proxy-authenticate", []string{"Basic realm=\"Proxy\""}, true},
		{"age", "age", []string{"3600"}, true},
		{"cache-control", "cache-control", []string{"no-cache"}, true},
		{"expires", "expires", []string{"Mon, 09 Dec 2025 12:00:00 GMT"}, true},
		{"last-modified", "last-modified", []string{"Mon, 09 Dec 2025 10:00:00 GMT"}, true},
		{"etag", "etag", []string{"\"abc123\""}, true},
		{"accept-ranges", "accept-ranges", []string{"bytes"}, true},
		{"content-range", "content-range", []string{"bytes 200-1023/146515"}, true},
		{"content-encoding", "content-encoding", []string{"gzip"}, true},
		{"content-language", "content-language", []string{"en-US"}, true},
		{"Content-Length mixed case", "Content-Length", []string{"100"}, true},
		{"CONTENT-LENGTH uppercase", "CONTENT-LENGTH", []string{"100"}, true},
		{"Set-Cookie mixed case", "Set-Cookie", []string{"session=abc123"}, true},
		{"SERVER uppercase", "SERVER", []string{"apache"}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			resourceHeaders := resource.NewResourceHeaders()
			resourceHeaders.Add(ctx, tc.headerName, tc.headerValues)

			headers := *resourceHeaders.Headers()
			_, exists := headers[tc.headerName]

			if tc.shouldBlock && exists {
				t.Fatalf("Expected header '%s' to be blocked, but it was added", tc.headerName)
			}
			if !tc.shouldBlock && !exists {
				t.Fatalf("Expected header '%s' to be allowed, but it was blocked", tc.headerName)
			}
		})
	}
}

func TestAddShouldAllowNonReservedHeaders(t *testing.T) {
	testCases := []struct {
		name         string
		headerName   string
		headerValues []string
	}{
		{"custom header", "X-Custom-Header", []string{"allowed"}},
		{"custom app header", "X-App-Version", []string{"1.2.3"}},
		{"multiple values", "X-Multi", []string{"value1", "value2", "value3"}},
		{"another custom header", "MyHeader", []string{"value1"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			resourceHeaders := resource.NewResourceHeaders()

			resourceHeaders.Add(ctx, tc.headerName, tc.headerValues)

			headers := *resourceHeaders.Headers()
			values, exists := headers[tc.headerName]

			if !exists {
				t.Fatalf("Expected header '%s' to be allowed, but it was blocked", tc.headerName)
			}
			if len(values) != len(tc.headerValues) {
				t.Fatalf("Expected %d values for header '%s', got %d", len(tc.headerValues), tc.headerName, len(values))
			}
			for i, expectedValue := range tc.headerValues {
				if values[i] != expectedValue {
					t.Fatalf("Expected value '%s' at index %d for header '%s', got '%s'", expectedValue, i, tc.headerName, values[i])
				}
			}
		})
	}
}
