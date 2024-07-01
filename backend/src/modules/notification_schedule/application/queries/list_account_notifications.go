package notification_schedule_query

import (
	notification_schedule_daos "weather_notification/src/modules/notification_schedule/application/daos"
	application_service "weather_notification/src/modules/shared/application"
	application_error "weather_notification/src/modules/shared/application/errors"
)

type ListAccountNotificationsImp struct {
	notificationDAO notification_schedule_daos.NotificationDAO
}

func NewListAccountNotifications(notificationDAO notification_schedule_daos.NotificationDAO) ListAccountNotifications {
	return ListAccountNotificationsImp{notificationDAO: notificationDAO}
}

func (service ListAccountNotificationsImp) Execute(input ListAccountNotificationsInputDTO) application_service.Output[ListAccountNotificationsOutputDTO] {
	notifications, err := service.notificationDAO.FindNotifications(input.AccountId)

	if err != nil {
		return application_service.Output[ListAccountNotificationsOutputDTO]{
			Error: &application_error.ApplicationError{
				Message: "could no query notifications",
				Err:     err,
				Name:    "CouldNotQueryNotifications",
			},
		}
	}

	return application_service.Output[ListAccountNotificationsOutputDTO]{
		Result: &ListAccountNotificationsOutputDTO{
			Notifications: notifications,
		},
		Error: nil,
	}
}

type ListAccountNotifications = application_service.ApplicationService[
	ListAccountNotificationsInputDTO,
	application_service.Output[ListAccountNotificationsOutputDTO],
]

type ListAccountNotificationsInputDTO struct {
	AccountId string
}

type ListAccountNotificationsOutputDTO struct {
	Notifications []notification_schedule_daos.NotificationDTO
}
