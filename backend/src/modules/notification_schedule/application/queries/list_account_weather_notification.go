package notification_schedule_query

import (
	"time"
	notification_schedule_daos "weather_notification/src/modules/notification_schedule/application/daos"
	application_service "weather_notification/src/modules/shared/application"
	application_error "weather_notification/src/modules/shared/application/errors"
)

type ListAccountWeatherNotificationServiceImp struct {
	notificationScheduleDAO notification_schedule_daos.ScheduledNotificationDAO
}

func (service ListAccountWeatherNotificationServiceImp) Execute(input ListAccountWeatherNotificationsServiceInputDTO) application_service.Output[[]ListAccountWeatherNotificationsServiceOutputDTO] {
	scheduledNotifications, err := service.notificationScheduleDAO.FindByAccountId(input.AccountId)

	if err != nil {
		return application_service.Output[[]ListAccountWeatherNotificationsServiceOutputDTO]{
			Error: &application_error.ApplicationError{
				Message: "Could not query for scheduled notifications",
				Err:     err,
				Name:    "ListAccountWeatherNotificationServiceError",
			},
			Result: nil,
		}
	}

	result := make([]ListAccountWeatherNotificationsServiceOutputDTO, 0, len(scheduledNotifications))

	for _, scheduledNotification := range scheduledNotifications {
		result = append(result, ListAccountWeatherNotificationsServiceOutputDTO{
			Id:             scheduledNotification.Id,
			ScheduledDate:  scheduledNotification.ScheduledDate,
			IntervalInDays: scheduledNotification.IntervalInDays,
			Hour:           scheduledNotification.Hour,
			City: ListAccountWeatherNotificationServiceCityOutputDTO{
				Id:        scheduledNotification.City.Id,
				Name:      scheduledNotification.City.Name,
				StateCode: scheduledNotification.City.StateCode,
				IsCoastal: scheduledNotification.IsCoastalCity,
			},
			Active: scheduledNotification.Active,
		})
	}

	return application_service.Output[[]ListAccountWeatherNotificationsServiceOutputDTO]{
		Error: nil, Result: &result,
	}
}

func NewListAccountWeatherNotificationService(notificationScheduleDAO notification_schedule_daos.ScheduledNotificationDAO) ListAccountWeatherNotificationsService {
	return ListAccountWeatherNotificationServiceImp{notificationScheduleDAO: notificationScheduleDAO}
}

type ListAccountWeatherNotificationsService = application_service.ApplicationService[
	ListAccountWeatherNotificationsServiceInputDTO,
	application_service.Output[[]ListAccountWeatherNotificationsServiceOutputDTO],
]

type ListAccountWeatherNotificationsServiceInputDTO struct {
	AccountId string
}

type ListAccountWeatherNotificationsServiceOutputDTO struct {
	Id             string
	ScheduledDate  time.Time
	IntervalInDays uint8
	Hour           uint8
	City           ListAccountWeatherNotificationServiceCityOutputDTO
	Active         bool
}

type ListAccountWeatherNotificationServiceCityOutputDTO struct {
	Id        string
	Name      string
	StateCode string
	IsCoastal bool
}
