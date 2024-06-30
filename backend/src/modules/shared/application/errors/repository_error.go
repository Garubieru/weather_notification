package application_error

import (
	"fmt"
)

func NewRepositoryError(name string) *ApplicationError {
	return &ApplicationError{
		Message: fmt.Sprintf("%s error", name),
		Name:    name,
		Err:     fmt.Errorf("%s error", name)}
}
