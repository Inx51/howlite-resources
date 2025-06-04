package errors

import "strings"

type aggregatedError struct {
	errors []error
}

func NewAggregatedError(errors []error) *aggregatedError {
	return &aggregatedError{
		errors: errors,
	}
}

func (aggregatedError *aggregatedError) Error() string {
	if len(aggregatedError.errors) == 0 {
		return "no errors"
	}

	errorMessages := make([]string, len(aggregatedError.errors))
	for i, err := range aggregatedError.errors {
		errorMessages[i] = err.Error()
	}

	return "aggregated error: " + strings.Join(errorMessages, "; ")
}
