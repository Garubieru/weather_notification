package infra_event

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	notification_schedule_events "weather_notification/src/modules/notification_schedule/application/events"
	event_broker "weather_notification/src/modules/notification_schedule/domain/event"
	infra_database "weather_notification/src/modules/shared/infra/database"

	"github.com/go-redis/redis/v8"
)

type RedisAccountNotificationBroker struct {
	redisClient redis.Client
	context     context.Context
	database    infra_database.Database
}

func (broker RedisAccountNotificationBroker) Emit(event event_broker.Event) error {
	eventPayload := reflect.ValueOf(event.Payload)

	switch eventPayload.Kind() {
	case reflect.Struct:
		eventPayload := event.Payload.(notification_schedule_events.AccountNotificationPayload)

		err := broker.database.Transaction(func(tx *sql.Tx) error {
			queryBuilder := broker.database.QueryBuilder("notification_schedule").
				SetColumns([]string{"scheduled_date"}).
				Where("account_id = ? AND id = ?")

			err := broker.database.Exec(queryBuilder.Update(),
				eventPayload.NextScheduledDate,
				eventPayload.AccountId,
				eventPayload.ScheduleId)

			if err != nil {
				return errors.New("could not save next schedule date")
			}

			queryBuilder = broker.database.QueryBuilder("notification").SetColumns(
				[]string{"id", "payload", "account_id", "seen"},
			)

			eventJson, err := json.Marshal(event)

			if err != nil {
				return errors.New("could not parse data")
			}

			err = broker.database.Exec(queryBuilder.Insert(1), event.Id, eventJson, eventPayload.AccountId, false)

			if err != nil {
				return errors.New("could not save notification")
			}

			channel := fmt.Sprintf("channel_account_%s", eventPayload.AccountId)

			output := broker.redisClient.Publish(broker.context, channel, eventJson)

			if output.Err() != nil {
				return output.Err()
			}

			return nil
		})

		return err
	default:
		return fmt.Errorf("unsupported type: %s", eventPayload.Kind().String())
	}
}

func (broker RedisAccountNotificationBroker) Subscribe(eventName string, handler func(message []byte) error) error {
	sub := broker.redisClient.Subscribe(broker.context, fmt.Sprintf("channel_account_%s", eventName))
	_, err := sub.Receive(broker.context)

	if err != nil {
		return err
	}

	defer sub.Close()

	channel := sub.Channel()

	for msg := range channel {
		fmt.Printf("Message to: %s, values: %s\n", eventName, string(msg.Payload))
		handler([]byte(msg.Payload))
	}

	return nil
}

func NewRedisAccountNotificationBroker(
	redisClient redis.Client,
	context context.Context,
	database infra_database.Database) RedisAccountNotificationBroker {
	return RedisAccountNotificationBroker{
		redisClient: redisClient,
		context:     context,
		database:    database,
	}
}
