package notification_schedule_services

import (
	"errors"
	"time"
	notification_schedule_daos "weather_notification/src/modules/notification_schedule/application/daos"
	notification_schedule_events "weather_notification/src/modules/notification_schedule/application/events"
	event_broker "weather_notification/src/modules/notification_schedule/domain/event"
)

type SendScheduledNotification struct {
	scheduledNotificationDAO notification_schedule_daos.ScheduledNotificationDAO
	eventBroker              event_broker.EventBroker
}

func (service SendScheduledNotification) Execute() error {
	notifications, err := service.scheduledNotificationDAO.FindByScheduledDate(time.Now())

	if err != nil {
		return errors.New("could not query for notifications")
	}

	if len(notifications) == 0 {
		return nil
	}

	if err := service.eventBroker.Emit(notification_schedule_events.NewScheduledNotificationsEvent(notifications)); err != nil {
		return errors.New("could not send notifications")
	}

	return nil
}

func NewSendScheduleNotification(
	scheduledNotificationDAO notification_schedule_daos.ScheduledNotificationDAO,
	eventBroker event_broker.EventBroker,
) SendScheduledNotification {
	return SendScheduledNotification{scheduledNotificationDAO: scheduledNotificationDAO, eventBroker: eventBroker}
}
