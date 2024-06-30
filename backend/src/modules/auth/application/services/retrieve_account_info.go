package auth_application_service

import (
	"errors"
	"weather_notification/src/modules/auth/domain/repositories"
	application_service "weather_notification/src/modules/shared/application"
	application_error "weather_notification/src/modules/shared/application/errors"
	"weather_notification/src/modules/shared/value_objects"
)

type RetrieveAccountInfoImp struct {
	accountRepository repositories.AccountRepository
}

func (service RetrieveAccountInfoImp) Execute(input RetrieveAccountInfoInputDTO) application_service.Output[RetrieveAccountInfoOutputDTO] {
	account, err := service.accountRepository.FindById(value_objects.RecoverID(input.AccountId))

	if err != nil {
		return application_service.Output[RetrieveAccountInfoOutputDTO]{
			Error:  application_error.NewRepositoryError("AccountRepository"),
			Result: nil,
		}
	}

	if account == nil {
		return application_service.Output[RetrieveAccountInfoOutputDTO]{
			Error: &application_error.ApplicationError{
				Message: "Account not found",
				Err:     errors.New("account not found"),
				Name:    "AccountNotFound",
			},
			Result: nil,
		}
	}

	return application_service.Output[RetrieveAccountInfoOutputDTO]{
		Result: &RetrieveAccountInfoOutputDTO{
			Id:       account.Id.Value,
			Name:     account.Name,
			Email:    account.Email.Value,
			Username: account.Username,
			Phone:    account.Phone,
		},
		Error: nil,
	}
}

type RetrieveAccountInfo = application_service.ApplicationService[
	RetrieveAccountInfoInputDTO,
	application_service.Output[RetrieveAccountInfoOutputDTO],
]

type RetrieveAccountInfoInputDTO struct {
	AccountId string
}

type RetrieveAccountInfoOutputDTO struct {
	Id       string
	Name     string
	Email    string
	Phone    string
	Username string
}

func NewRetrieveAccountInfo(accountRepository repositories.AccountRepository) RetrieveAccountInfo {
	return RetrieveAccountInfoImp{accountRepository: accountRepository}
}
