package resource

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
