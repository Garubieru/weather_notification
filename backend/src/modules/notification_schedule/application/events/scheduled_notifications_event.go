package notification_schedule_events

import (
	notification_schedule_daos "weather_notification/src/modules/notification_schedule/application/daos"
	event_broker "weather_notification/src/modules/notification_schedule/domain/event"
	"weather_notification/src/modules/shared/value_objects"
)

func NewScheduledNotificationsEvent(payload []notification_schedule_daos.ScheduledNotificationDTO) event_broker.Event {
	return event_broker.Event{
		Id:      value_objects.NewID().Value,
		Name:    "SendScheduledNotifications",
		Payload: payload,
	}
}
