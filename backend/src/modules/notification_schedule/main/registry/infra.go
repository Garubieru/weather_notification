package notification_registry_pattern

import (
	"context"
	"weather_notification/src/modules/notification_schedule/application/event_handlers"
	infra_event "weather_notification/src/modules/notification_schedule/infra/event"
	infra_gateways "weather_notification/src/modules/notification_schedule/infra/gateways"
	registry "weather_notification/src/modules/shared/infra"
	infra_database "weather_notification/src/modules/shared/infra/database"

	"github.com/go-redis/redis/v8"
)

func RegisterInfra(ctx context.Context) {
	registry := registry.GetRegistryInstance()

	redisClient := *redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	registry.Register(InfraKeys.RedisClient, redisClient)

	kafkaBroker := infra_event.NewKafkaEventBroker(
		"schedule_weather_notification",
		0,
		"localhost:9092",
		context.WithoutCancel(ctx),
		registry.Inject("Database").(infra_database.Database),
	)

	redisAccountBroker := infra_event.NewRedisAccountNotificationBroker(
		redisClient,
		context.WithoutCancel(ctx),
	)

	sendNotificationHandler := event_handlers.NewSendAccountWeatherNotification(
		infra_gateways.CPTECWeatherGateway{},
		redisAccountBroker,
	)

	kafkaBroker.Subscribe("schedule_weather_notification", sendNotificationHandler.Handle)

	kafkaBroker.Listen()

	registry.Register(InfraKeys.KafkaEventBroker, kafkaBroker)
	registry.Register(InfraKeys.RedisEventBroker, redisAccountBroker)
}

type infraKeys struct {
	RedisClient      string
	RedisEventBroker string
	KafkaEventBroker string
}

var InfraKeys = infraKeys{
	RedisClient:      "RedisClient",
	RedisEventBroker: "RedisEventBroker",
	KafkaEventBroker: "KafkaEventBroker",
}
