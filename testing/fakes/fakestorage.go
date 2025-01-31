package fakes

import (
	"bytes"
	"io"

	"github.com/inx51/howlite/resources/resource"
)

type FakeStorage struct {
	storage map[resource.ResourceIdentifier]*memoryWriteCloser
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
		storage: make(map[resource.ResourceIdentifier]*memoryWriteCloser),
	}
}

func (storage *FakeStorage) RemoveResource(resourceIdentifier *resource.ResourceIdentifier) error {
	delete(storage.storage, *resourceIdentifier)
	return nil
}

func (storage *FakeStorage) NewResourceWriter(resourceIdentifier *resource.ResourceIdentifier) (io.WriteCloser, error) {
	buffer := &memoryWriteCloser{buffer: new(bytes.Buffer)}
	storage.storage[*resourceIdentifier] = buffer
	return buffer, nil
}

func (storage *FakeStorage) ResourceExists(resourceIdentifier *resource.ResourceIdentifier) (bool, error) {
	_, ok := storage.storage[*resourceIdentifier]
	return ok, nil
}

func (storage *FakeStorage) GetResource(resourceIdentifier *resource.ResourceIdentifier) (io.ReadCloser, error) {
	return storage.storage[*resourceIdentifier], nil
}
