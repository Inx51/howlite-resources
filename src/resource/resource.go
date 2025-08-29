package resource

import (
	"io"
)

type Resource struct {
	Identifier *ResourceIdentifier
	Headers    *ResourceHeaders
	Body       *io.ReadCloser
}

func NewResource(identifier *ResourceIdentifier, body *io.ReadCloser) *Resource {
	return &Resource{
		Identifier: identifier,
		Body:       body,
		Headers:    NewResourceHeaders(),
	}
}

func (resource *Resource) Write(writer io.WriteCloser) error {
	if err := resource.Headers.writeHeaders(writer); err != nil {
		return err
	}

	buff := make([]byte, 1024)
	readCloser := io.NopCloser(*resource.Body)
	_, err := io.CopyBuffer(writer, readCloser, buff)
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	return nil
}

func LoadResource(resourceIdentifier *ResourceIdentifier, reader io.ReadCloser) (*Resource, error) {
	resourceHeaders := NewResourceHeaders()
	err := resourceHeaders.LoadHeaders(reader)
	if err != nil {
		return nil, err
	}

	resource := &Resource{
		Identifier: resourceIdentifier,
		Headers:    resourceHeaders,
		Body:       &reader,
	}

	return resource, err
}
