package notification_schedule_services

import (
	"errors"
	notification_schedule_repositories "weather_notification/src/modules/notification_schedule/domain/repositories"
	application_service "weather_notification/src/modules/shared/application"
	application_error "weather_notification/src/modules/shared/application/errors"
	"weather_notification/src/modules/shared/value_objects"
)

type DeactivateWeatherNotificationScheduleServiceImp struct {
	accountRepository notification_schedule_repositories.AccountRepository
}

func (service DeactivateWeatherNotificationScheduleServiceImp) Execute(input DeactivateWeatherNotificationScheduleInputDTO) application_service.Output[DeactivateWeatherNotificationScheduleOutputDTO] {
	account, err := service.accountRepository.FindById(value_objects.RecoverID(input.AccountId))

	if err != nil {
		return application_service.Output[DeactivateWeatherNotificationScheduleOutputDTO]{
			Error:  application_error.NewRepositoryError("AccountRepository"),
			Result: nil,
		}
	}

	if account == nil {
		return application_service.Output[DeactivateWeatherNotificationScheduleOutputDTO]{
			Error:  &application_error.ApplicationError{Message: "account not found", Err: errors.New("account not found"), Name: "AccountNotFound"},
			Result: nil,
		}
	}

	deactivateError := account.DeactivateSchedule(input.ScheduleId)

	if deactivateError != nil {
		return application_service.Output[DeactivateWeatherNotificationScheduleOutputDTO]{
			Error:  &application_error.ApplicationError{Message: deactivateError.Error(), Err: deactivateError, Name: "DeactivateScheduleError"},
			Result: nil,
		}
	}

	if err := service.accountRepository.Save(account); err != nil {
		return application_service.Output[DeactivateWeatherNotificationScheduleOutputDTO]{
			Error:  application_error.NewRepositoryError("AccountRepository"),
			Result: nil,
		}
	}

	return application_service.Output[DeactivateWeatherNotificationScheduleOutputDTO]{
		Error:  nil,
		Result: &DeactivateWeatherNotificationScheduleOutputDTO{ScheduleId: input.ScheduleId},
	}
}

type DeactivateWeatherNotificationScheduleService = application_service.ApplicationService[DeactivateWeatherNotificationScheduleInputDTO, application_service.Output[DeactivateWeatherNotificationScheduleOutputDTO]]

type DeactivateWeatherNotificationScheduleInputDTO struct {
	ScheduleId string
	AccountId  string
}

type DeactivateWeatherNotificationScheduleOutputDTO struct {
	ScheduleId string
}

func NewDeactivateWeatherNotificationScheduleService(accountRepository notification_schedule_repositories.AccountRepository) DeactivateWeatherNotificationScheduleService {
	return DeactivateWeatherNotificationScheduleServiceImp{accountRepository: accountRepository}
}
