package auth_application_service

import (
	"fmt"
	"weather_notification/src/modules/auth/domain/repositories"
	application_service "weather_notification/src/modules/shared/application"
	application_error "weather_notification/src/modules/shared/application/errors"
	"weather_notification/src/modules/shared/value_objects"
)

type AuthenticateSessionServiceImp struct {
	sessionRepository repositories.SessionRepository
}

func (service AuthenticateSessionServiceImp) Execute(input AuthenticateSessionInput) application_service.Output[AuthenticateSessionOutput] {
	if len(input.SessionId) == 0 {
		return application_service.Output[AuthenticateSessionOutput]{
			Result: &AuthenticateSessionOutput{Authenticated: false},
			Error:  nil,
		}
	}

	session, err := service.sessionRepository.FindById(value_objects.RecoverID(input.SessionId))

	if err != nil {
		return application_service.Output[AuthenticateSessionOutput]{
			Error: application_error.NewRepositoryError("SessionRepository"),
		}
	}

	if session == nil || session.IsExpired() {
		if err := service.sessionRepository.Delete(session); err != nil {
			fmt.Printf("Could not delete session: %s\n", session.Id.Value)
		}

		return service.buildNotAuthenticateOutput()
	}

	return application_service.Output[AuthenticateSessionOutput]{Result: &AuthenticateSessionOutput{Authenticated: true, AccountId: &session.AccountId}}
}

func (service AuthenticateSessionServiceImp) buildNotAuthenticateOutput() application_service.Output[AuthenticateSessionOutput] {
	return application_service.Output[AuthenticateSessionOutput]{
		Result: &AuthenticateSessionOutput{Authenticated: false, AccountId: nil},
		Error:  nil,
	}
}

type AuthenticateSessionService = application_service.ApplicationService[AuthenticateSessionInput, application_service.Output[AuthenticateSessionOutput]]

type AuthenticateSessionInput struct {
	SessionId string
}

type AuthenticateSessionOutput struct {
	Authenticated bool
	AccountId     *value_objects.ID
}

func NewAuthenticateSessionService(sessionRepository repositories.SessionRepository) AuthenticateSessionService {
	return &AuthenticateSessionServiceImp{sessionRepository: sessionRepository}
}
