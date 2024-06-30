package notification_schedule_events

import (
	"time"
	"weather_notification/src/modules/notification_schedule/application/gateways"
	event_broker "weather_notification/src/modules/notification_schedule/domain/event"
	"weather_notification/src/modules/shared/value_objects"
)

func NewAccountNotificationEvent(payload AccountNotificationPayload) event_broker.Event {
	return event_broker.Event{
		Id:      value_objects.NewID().Value,
		Name:    "AccountNotificationEvent",
		Payload: payload,
	}
}

type AccountNotificationPayload struct {
	AccountId         string                 `json:"accountId"`
	ScheduleId        string                 `json:"scheduleId"`
	CityName          string                 `json:"cityName"`
	NextScheduledDate time.Time              `json:"nextScheduledDate"`
	CityStateCode     string                 `json:"cityStateCode"`
	Prediction        gateways.PredictionDTO `json:"prediction"`
}
