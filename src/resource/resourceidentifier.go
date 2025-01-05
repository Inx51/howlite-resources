package resource

type ResourceIdentifier struct {
	Value *string
}

func NewIdentifier(identifier *string) ResourceIdentifier {
	return ResourceIdentifier{Value: identifier}
}
