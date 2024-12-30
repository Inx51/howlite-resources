package errors

import (
	"fmt"
)

type NotFoundError struct {
	Identifier *string
}

// Error implements error.
func (r *NotFoundError) Error() string {
	return fmt.Sprintf("Could not find resource with identifier %s", *r.Identifier)
}

// var (
// 	ErrResourceNotFound = errNotFound()
// )

// type ResourceNotFound struct {
// 	Identifier string
// }

// // Error implements error.
// func (r ResourceNotFound) Error() string {
// 	return fmt.Sprintf("Could not find resource with identifier %s", r.Identifier)
// }

// func errNotFound(r ResourceNotFound) error {
// 	return errors.New()
// }
