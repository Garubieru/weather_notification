package application_error

import "fmt"

type ApplicationError struct {
	Message string
	Err     error
	Name    string
}

func (e *ApplicationError) Error() string {
	return fmt.Sprintf("%s failed: %v", e.Message, e.Err)
}
