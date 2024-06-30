package notification_schedule_daos

import (
	event_broker "weather_notification/src/modules/notification_schedule/domain/event"
)

type NotificationDAO interface {
	FindNotifications(accountId string) ([]NotificationDTO, error)
}

type NotificationDTO = event_broker.Event
