package s3

type UndefinedStrategyError struct {
	strategy string
}

func NewUndefinedStrategyError(strategy string) *UndefinedStrategyError {
	return &UndefinedStrategyError{
		strategy: strategy,
	}
}

func (e *UndefinedStrategyError) Error() string {
	return "Undefined strategy for S3 storage: " + e.strategy
}
