package notification_schedule_services

import (
	"errors"
	"time"
	notification_schedule_daos "weather_notification/src/modules/notification_schedule/application/daos"
	notification_domain_entities "weather_notification/src/modules/notification_schedule/domain/entities"
	notification_schedule_repositories "weather_notification/src/modules/notification_schedule/domain/repositories"
	application_service "weather_notification/src/modules/shared/application"
	application_error "weather_notification/src/modules/shared/application/errors"
	"weather_notification/src/modules/shared/value_objects"
)

type ScheduleNotificationApplicationServiceImp struct {
	accountRepository notification_schedule_repositories.AccountRepository
	cityDAO           notification_schedule_daos.CityDAO
}

func (service ScheduleNotificationApplicationServiceImp) Execute(input ScheduleNotificationInputDTO) application_service.Output[ScheduleNotificationOutputDTO] {
	account, err := service.accountRepository.FindById(value_objects.RecoverID(input.AccountId))

	if err != nil {
		return application_service.Output[ScheduleNotificationOutputDTO]{
			Error:  application_error.NewRepositoryError("AccountRepository"),
			Result: nil,
		}
	}

	if account == nil {
		return application_service.Output[ScheduleNotificationOutputDTO]{
			Error: &application_error.ApplicationError{
				Message: "Account not found",
				Err:     errors.New("account not found"),
				Name:    "AccountNotFound",
			},
			Result: nil,
		}
	}

	city, err := service.cityDAO.FindById(input.CityId)

	if err != nil {
		return application_service.Output[ScheduleNotificationOutputDTO]{
			Error: &application_error.ApplicationError{
				Message: "Error on retrieving city data",
				Err:     errors.New("error on retrieving city data"),
				Name:    "CityQueryError",
			},
			Result: nil,
		}
	}

	if city == nil {
		return application_service.Output[ScheduleNotificationOutputDTO]{
			Error: &application_error.ApplicationError{
				Message: "City not found",
				Err:     errors.New("city not found"),
				Name:    "CityNotFound",
			},
			Result: nil,
		}
	}

	scheduleError := account.ScheduleWeatherNotification(notification_domain_entities.ScheduleWeatherNotificationInput{
		CityId:         city.Id,
		IsCityCoastal:  input.IsCityCoastal,
		IntervalInDays: input.IntervalInDays,
		Method:         input.Method,
		Hour:           input.Hour,
	})

	if scheduleError != nil {
		return application_service.Output[ScheduleNotificationOutputDTO]{
			Error: &application_error.ApplicationError{
				Message: scheduleError.Error(),
				Err:     err,
				Name:    "ScheduleNotificationError",
			},
			Result: nil,
		}
	}

	if err := service.accountRepository.Save(account); err != nil {
		return application_service.Output[ScheduleNotificationOutputDTO]{
			Error: &application_error.ApplicationError{
				Message: err.Error(),
				Err:     err,
				Name:    "SaveAccountRepositoryError",
			},
			Result: nil,
		}
	}

	newSchedule := account.ScheduledNotifications.GetNewItems()[0]

	return application_service.Output[ScheduleNotificationOutputDTO]{
		Error:  nil,
		Result: &ScheduleNotificationOutputDTO{ScheduleTime: newSchedule.ScheduledDate, Id: newSchedule.Id.Value},
	}
}

type ScheduleNotificationApplicationService = application_service.ApplicationService[ScheduleNotificationInputDTO, application_service.Output[ScheduleNotificationOutputDTO]]

type ScheduleNotificationInputDTO struct {
	AccountId      string
	CityId         string
	IsCityCoastal  bool
	IntervalInDays uint8
	Hour           uint8
	Method         string
}

type ScheduleNotificationOutputDTO struct {
	Id           string
	ScheduleTime time.Time
}

func NewScheduleWeatherNotification(
	accountRepository notification_schedule_repositories.AccountRepository,
	cityDao notification_schedule_daos.CityDAO,
) ScheduleNotificationApplicationService {
	return ScheduleNotificationApplicationServiceImp{
		accountRepository: accountRepository,
		cityDAO:           cityDao,
	}
}
