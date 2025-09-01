package resource_test

import (
	"testing"

	"github.com/inx51/howlite-resources/resource"
)

func TestAdd_Should_add_headers(t *testing.T) {
	resourceHeaders := resource.NewResourceHeaders()
	resourceHeaders.Add("A", []string{"1", "2"})
	resourceHeaders.Add("B", []string{"3"})
	headers := resourceHeaders.Headers()
	if len(*headers) != 2 {
		t.Fatal("Headers() length mismatch, got", len(*headers), "want 2")
	}
}

func TestAdd_Should_have_values(t *testing.T) {
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

// func TestWriteHeaders_Empty(t *testing.T) {
// 	rh := resourceHeaders.NewResourceHeaders()
// 	var buf bytes.Buffer
// 	err := rh.writeHeaders(&buf)
// 	if err != nil {
// 		t.Fatalf("writeHeaders failed: %v", err)
// 	}
// 	b := buf.Bytes()
// 	if len(b) != 8 || !bytes.Equal(b, make([]byte, 8)) {
// 		t.Errorf("writeHeaders for empty headers should write 8 zero bytes, got: %v", b)
// 	}
// }

// func TestWriteHeaders_NonEmpty(t *testing.T) {
// 	rh := NewResourceHeaders()
// 	rh.Add("foo", []string{"bar"})
// 	var buf bytes.Buffer
// 	err := rh.writeHeaders(&buf)
// 	if err != nil {
// 		t.Fatalf("writeHeaders failed: %v", err)
// 	}
// 	b := buf.Bytes()
// 	if len(b) < 8 {
// 		t.Fatalf("writeHeaders output too short: %d", len(b))
// 	}
// 	// First 8 bytes: length
// 	msgLen := binaryLittleEndianUint64(b[:8])
// 	if int(msgLen) != len(b[8:]) {
// 		t.Errorf("msgpack length mismatch: header says %d, actual %d", msgLen, len(b[8:]))
// 	}
// 	// Try to decode
// 	var decoded map[string][]string
// 	err = msgpack.Unmarshal(b[8:], &decoded)
// 	if err != nil {
// 		t.Errorf("msgpack decode failed: %v", err)
// 	}
// 	if !reflect.DeepEqual(decoded, rh.headers) {
// 		t.Errorf("decoded headers mismatch: got %v, want %v", decoded, rh.headers)
// 	}
// }

func TestLoadHeadersAndGetHeadersFromStream(t *testing.T) {
	// // Prepare headers
	// orig := map[string][]string{"a": {"b", "c"}}
	// packed, err := msgpack.Marshal(orig)
	// if err != nil {
	// 	t.Fatalf("msgpack marshal failed: %v", err)
	// }
	// var buf bytes.Buffer
	// // Write length
	// l := make([]byte, 8)
	// putBinaryLittleEndianUint64(l, uint64(len(packed)))
	// buf.Write(l)
	// buf.Write(packed)
	// // Use LoadHeaders
	// resourceHeaders := resource.NewResourceHeaders()
	// resourceHeadersReaderCloser := io.NopCloser(bytes.NewReader(buf.Bytes()))
	// err = resourceHeaders.LoadHeaders(resourceHeadersReaderCloser)
	// if err != nil {
	// 	t.Fatalf("LoadHeaders failed: %v", err)
	// }
	// if !reflect.DeepEqual(resourceHeaders.Headers(), orig) {
	// 	t.Errorf("LoadHeaders did not load correctly: got %v, want %v", resourceHeaders.Headers(), orig)
	// }
}

func TestLoadHeaders_Empty(t *testing.T) {
	// // 8 zero bytes
	// buf := bytes.NewBuffer(make([]byte, 8))
	// rh := NewResourceHeaders()
	// rc := io.NopCloser(buf)
	// err := rh.LoadHeaders(rc)
	// if err != nil {
	// 	t.Fatalf("LoadHeaders failed for empty: %v", err)
	// }
	// if len(rh.headers) != 0 {
	// 	t.Errorf("LoadHeaders for empty should result in empty map, got: %v", rh.headers)
	// }
}

func TestGetHeadersFromStream_Error(t *testing.T) {
	// // Not enough bytes for length
	// rc := io.NopCloser(bytes.NewReader([]byte{1, 2, 3}))
	// _, err := getHeadersFromStream(rc)
	// if err == nil {
	// 	t.Error("expected error for short stream, got nil")
	// }
}
