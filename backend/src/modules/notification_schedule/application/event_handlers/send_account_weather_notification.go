package event_handlers

import (
	"encoding/json"
	"errors"
	notification_schedule_daos "weather_notification/src/modules/notification_schedule/application/daos"
	notification_schedule_events "weather_notification/src/modules/notification_schedule/application/events"
	"weather_notification/src/modules/notification_schedule/application/gateways"
	event_broker "weather_notification/src/modules/notification_schedule/domain/event"
)

type SendAccountWeatherNotification struct {
	weatherGateway gateways.WeatherGateway
	eventBroker    event_broker.EventBroker
}

func (handler SendAccountWeatherNotification) Handle(message []byte) error {
	var messageData notification_schedule_daos.ScheduledNotificationDTO

	err := json.Unmarshal(message, &messageData)

	if err != nil {
		return ErrEventHandlerParseData
	}

	predictions, err := handler.weatherGateway.GetPrediction(4, messageData.City.ExternalId, messageData.IsCoastalCity)

	if err != nil {
		return ErrEventHandlerUnavailable
	}

	event := notification_schedule_events.NewAccountNotificationEvent(notification_schedule_events.AccountNotificationPayload{
		AccountId:     messageData.AccountId,
		CityName:      messageData.City.Name,
		CityStateCode: messageData.City.StateCode,
		Predictions:   predictions,
	})

	if emitError := handler.eventBroker.Emit(event); emitError != nil {
		return ErrEventHandlerUnavailable
	}

	return nil
}

var (
	ErrEventHandlerUnavailable = errors.New("event handler unavailable")
	ErrEventHandlerParseData   = errors.New("could not parse data")
)

func NewSendAccountWeatherNotification(
	weatherGateway gateways.WeatherGateway,
	eventBroker event_broker.EventBroker,
) SendAccountWeatherNotification {
	return SendAccountWeatherNotification{weatherGateway: weatherGateway, eventBroker: eventBroker}
}
