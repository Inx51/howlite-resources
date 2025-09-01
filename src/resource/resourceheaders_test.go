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

func TestAddShouldAddHeaders(t *testing.T) {
	resourceHeaders := resource.NewResourceHeaders()
	resourceHeaders.Add("A", []string{"1", "2"})
	resourceHeaders.Add("B", []string{"3"})

	headers := *resourceHeaders.Headers()

	if len(headers) != 2 {
		t.Fatal("Headers() length mismatch, got", len(headers), "want 2")
	}
}

func TestAddShouldHaveValues(t *testing.T) {
	resourceHeaders := resource.NewResourceHeaders()
	resourceHeaders.Add("A", []string{"1", "2"})
	resourceHeaders.Add("B", []string{"3"})

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
