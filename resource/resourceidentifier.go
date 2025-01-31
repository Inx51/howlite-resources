package resource

type ResourceIdentifier struct {
	Value *string
}

func NewResourceIdentifier(value *string) *ResourceIdentifier {
	return &ResourceIdentifier{
		Value: value,
	}
}
