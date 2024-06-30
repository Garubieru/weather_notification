package notification_schedule_services

import (
	"errors"
	notification_schedule_repositories "weather_notification/src/modules/notification_schedule/domain/repositories"
	application_service "weather_notification/src/modules/shared/application"
	application_error "weather_notification/src/modules/shared/application/errors"
	"weather_notification/src/modules/shared/value_objects"
)

type ActivateWeatherNotificationScheduleServiceImp struct {
	accountRepository notification_schedule_repositories.AccountRepository
}

func (service ActivateWeatherNotificationScheduleServiceImp) Execute(input ActivateWeatherNotificationScheduleInputDTO) application_service.Output[ActivateWeatherNotificationScheduleOutputDTO] {
	account, err := service.accountRepository.FindById(value_objects.RecoverID(input.AccountId))

	if err != nil {
		return application_service.Output[ActivateWeatherNotificationScheduleOutputDTO]{
			Error:  application_error.NewRepositoryError("AccountRepository"),
			Result: nil,
		}
	}

	if account == nil {
		return application_service.Output[ActivateWeatherNotificationScheduleOutputDTO]{
			Error:  &application_error.ApplicationError{Message: "account not found", Err: errors.New("account not found"), Name: "AccountNotFound"},
			Result: nil,
		}
}

	activateError := account.ActivateSchedule(input.ScheduleId)

	if activateError != nil {
		return application_service.Output[ActivateWeatherNotificationScheduleOutputDTO]{
			Error:  &application_error.ApplicationError{Message: activateError.Error(), Err: activateError, Name: "ActivateScheduleError"},
			Result: nil,
		}
	}

	if err := service.accountRepository.Save(account); err != nil {
		return application_service.Output[ActivateWeatherNotificationScheduleOutputDTO]{
			Error:  application_error.NewRepositoryError("AccountRepository"),
			Result: nil,
		}
	}

	return application_service.Output[ActivateWeatherNotificationScheduleOutputDTO]{
		Error:  nil,
		Result: &ActivateWeatherNotificationScheduleOutputDTO{ScheduleId: input.ScheduleId},
	}
}

type ActivateWeatherNotificationScheduleService = application_service.ApplicationService[ActivateWeatherNotificationScheduleInputDTO, application_service.Output[ActivateWeatherNotificationScheduleOutputDTO]]

type ActivateWeatherNotificationScheduleInputDTO struct {
	ScheduleId string
	AccountId  string
}

type ActivateWeatherNotificationScheduleOutputDTO struct {
	ScheduleId string
}

func NewActivateWeatherNotificationScheduleService(accountRepository notification_schedule_repositories.AccountRepository) ActivateWeatherNotificationScheduleService {
	return ActivateWeatherNotificationScheduleServiceImp{accountRepository: accountRepository}
}
