package application_service

import application_error "weather_notification/src/modules/shared/application/errors"

type ApplicationService[Input any, Output any] interface {
	Execute(input Input) Output
}

type Output[Result any] struct {
	Error  *application_error.ApplicationError
	Result *Result
}
