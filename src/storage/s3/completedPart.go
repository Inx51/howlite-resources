package s3

type CompletedPart struct {
	ETag       *string
	PartNumber *int32
}

func NewCompletedPart(etag *string, partNumber *int32) CompletedPart {
	return CompletedPart{
		ETag:       etag,
		PartNumber: partNumber,
	}
}
