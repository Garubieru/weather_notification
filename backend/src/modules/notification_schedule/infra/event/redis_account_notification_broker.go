package infra_event

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	notification_schedule_events "weather_notification/src/modules/notification_schedule/application/events"
	event_broker "weather_notification/src/modules/notification_schedule/domain/event"

	"github.com/go-redis/redis/v8"
)

type RedisAccountNotificationBroker struct {
	redisClient redis.Client
	context     context.Context
}

func (broker RedisAccountNotificationBroker) Emit(event event_broker.Event) error {
	eventPayload := reflect.ValueOf(event.Payload)

	switch eventPayload.Kind() {
	case reflect.Struct:
		eventPayload := event.Payload.(notification_schedule_events.AccountNotificationPayload)

		jsonData, err := json.Marshal(event)

		if err != nil {
			return errors.New("could not parse data")
		}

		channel := fmt.Sprintf("channel_account_%s", eventPayload.AccountId)

		broker.redisClient.Publish(broker.context, channel, jsonData)
	default:
		return fmt.Errorf("unsupported type: %s", eventPayload.Kind().String())
	}

	return nil
}

func (broker RedisAccountNotificationBroker) Subscribe(eventName string, handler func(message []byte) error) error {
	sub := broker.redisClient.Subscribe(broker.context, fmt.Sprintf("channel_account_%s", eventName))
	_, err := sub.Receive(broker.context)

	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("slaaaaa", eventName)

	defer sub.Close()

	channel := sub.Channel()

	for msg := range channel {
		fmt.Println(msg)
		handler([]byte(msg.Payload))
	}

	return nil
}

func NewRedisAccountNotificationBroker(redisClient redis.Client, context context.Context) RedisAccountNotificationBroker {
	return RedisAccountNotificationBroker{redisClient: redisClient, context: context}
}
