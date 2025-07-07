package fakes

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"

	"github.com/inx51/howlite/resources/resource"
	"github.com/vmihailenco/msgpack/v5"
)

type FakeStorage struct {
	storage map[string]*memoryWriteCloser
}

type memoryWriteCloser struct {
	buffer *bytes.Buffer
}

func (mw *memoryWriteCloser) Read(p []byte) (n int, err error) {
	return mw.buffer.Read(p)
}

func (mw *memoryWriteCloser) Write(p []byte) (n int, err error) {
	return mw.buffer.Write(p)
}

func (mw *memoryWriteCloser) Close() error {
	// No-op, because we don't need to free any resources
	return nil
}

func NewStorage() *FakeStorage {
	return &FakeStorage{
		storage: make(map[string]*memoryWriteCloser),
	}
}

func (storage *FakeStorage) AddTestResource(identifier string, headers map[string][]string, body []byte) error {

	//TODO: Could this be simplified?
	resIdentifier := resource.NewResourceIdentifier(&identifier)
	memWriteCloser, err := storage.NewResourceWriterContext(context.Background(), resIdentifier)
	if err != nil {
		panic(err)
	}
	defer memWriteCloser.Close()

	if headers == nil {
		headers = make(map[string][]string)
	}
	if body == nil {
		body = []byte{}
	}

	if len(headers) == 0 {
		if _, err := memWriteCloser.Write(make([]byte, 8)); err != nil {
			panic(err)
		}
	} else {
		msgPackedHeaders, err := msgpack.Marshal(headers)
		if err != nil {
			panic(err)
		}
		headersMsgPackLength := len(msgPackedHeaders)
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(headersMsgPackLength))
		if _, err := memWriteCloser.Write(b); err != nil {
			panic(err)
		}
		if _, err := memWriteCloser.Write(msgPackedHeaders); err != nil {
			panic(err)
		}
	}

	if _, err := memWriteCloser.Write(body); err != nil {
		panic(err)
	}

	return nil
}

func (storage *FakeStorage) GetName() string {
	return "FakeStorage"
}

func (storage *FakeStorage) RemoveResourceContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) error {
	delete(storage.storage, *resourceIdentifier.Value)
	return nil
}

func (storage *FakeStorage) NewResourceWriterContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (io.WriteCloser, error) {
	buffer := &memoryWriteCloser{buffer: new(bytes.Buffer)}
	storage.storage[*resourceIdentifier.Value] = buffer
	return buffer, nil
}

func (storage *FakeStorage) ResourceExistsContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (bool, error) {
	_, ok := storage.storage[*resourceIdentifier.Value]
	return ok, nil
}

func (storage *FakeStorage) GetResourceContext(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (io.ReadCloser, error) {
	return storage.storage[*resourceIdentifier.Value], nil
}
