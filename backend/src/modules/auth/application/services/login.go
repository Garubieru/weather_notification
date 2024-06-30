package auth_application_service

import (
	"errors"
	"weather_notification/src/modules/auth/domain/entities"
	"weather_notification/src/modules/auth/domain/repositories"
	services "weather_notification/src/modules/auth/domain/services"
	application_service "weather_notification/src/modules/shared/application"
	application_error "weather_notification/src/modules/shared/application/errors"
)

type LoginServiceImp struct {
	accountRepository repositories.AccountRepository
	sessionRepository repositories.SessionRepository
	cryptoService     services.CryptoService
}

func (service *LoginServiceImp) Execute(input LoginServiceInput) application_service.Output[LoginServiceOutput] {
	account, err := service.accountRepository.FindByUsername(input.Username)

	if err != nil {
		return application_service.Output[LoginServiceOutput]{Error: application_error.NewRepositoryError("AccountRepository"), Result: nil}
	}

	if account == nil {
		return *service.invalidUsernameOrPassword()
	}

	account.Authenticate(service.cryptoService, input.Password)

	if !account.IsAuthenticated {
		return *service.invalidUsernameOrPassword()
	}

	session := entities.NewSession(entities.SessionCreateCommand{AccountId: account.Id})

	if err := service.sessionRepository.Save(&session); err != nil {
		return application_service.Output[LoginServiceOutput]{Error: application_error.NewRepositoryError("SessionRepository"), Result: nil}
	}

	return application_service.Output[LoginServiceOutput]{Result: &LoginServiceOutput{SessionId: session.Id.Value}, Error: nil}
}

func (service *LoginServiceImp) invalidUsernameOrPassword() *application_service.Output[LoginServiceOutput] {
	return &application_service.Output[LoginServiceOutput]{Error: &application_error.ApplicationError{
		Message: "Invalid username or password",
		Err:     errors.New("invalid username or password"),
		Name:    "InvalidUsernameOrPassword",
	}, Result: nil}
}

type LoginService = application_service.ApplicationService[LoginServiceInput, application_service.Output[LoginServiceOutput]]

type LoginServiceInput struct {
	Username string
	Password string
}

type LoginServiceOutput struct {
	SessionId string
}

func NewLoginService(
	accountRepository repositories.AccountRepository,
	sessionRepository repositories.SessionRepository,
	cryptoSevice services.CryptoService) LoginService {
	return &LoginServiceImp{accountRepository: accountRepository, sessionRepository: sessionRepository, cryptoService: cryptoSevice}
}
