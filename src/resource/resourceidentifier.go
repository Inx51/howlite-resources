package resource

import (
	"crypto/md5"
	"encoding/base64"
)

type ResourceIdentifier struct {
	identifier string
}

func NewResourceIdentifier(identifier string) *ResourceIdentifier {
	return &ResourceIdentifier{
		identifier: identifier,
	}
}

func (resourceIdentifier *ResourceIdentifier) Identifier() string {
	return resourceIdentifier.identifier
}

func (resourceIdentifier *ResourceIdentifier) ToUniqueFilename() string {
	var encBytes = md5.Sum([]byte(resourceIdentifier.Identifier()))
	return base64.URLEncoding.EncodeToString(encBytes[:]) + ".bin"
}
