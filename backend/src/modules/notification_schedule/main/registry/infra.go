package notification_registry_pattern

import (
	"context"
	"log"
	"os"
	"strconv"
	"weather_notification/src/modules/notification_schedule/application/event_handlers"
	infra_event "weather_notification/src/modules/notification_schedule/infra/event"
	infra_gateways "weather_notification/src/modules/notification_schedule/infra/gateways"
	registry "weather_notification/src/modules/shared/infra"
	infra_database "weather_notification/src/modules/shared/infra/database"

	"github.com/go-redis/redis/v8"
)

func RegisterInfra(ctx context.Context) {
	registry := registry.GetRegistryInstance()

	redisHost := os.Getenv("REDIS_ADDR")

	redisClient := *redis.NewClient(&redis.Options{
		Addr: redisHost,
	})

	registry.Register(InfraKeys.RedisClient, redisClient)

	kafkaTopic := os.Getenv("KAKFA_TOPIC")
	kafkaHost := os.Getenv("KAFKA_HOST")
	kafkaPartition, err := strconv.ParseUint(os.Getenv("KAKFA_PARTITION"), 10, 8)

	if err != nil {
		log.Fatal("Provide a valid partition")
	}

	kafkaBroker := infra_event.NewKafkaEventBroker(infra_event.KafkaEventBrokerConfig{
		Topic:     kafkaTopic,
		Partition: int(kafkaPartition),
		Host:      kafkaHost,
		Context:   context.WithoutCancel(ctx),
		Database:  registry.Inject("Database").(infra_database.Database),
	},
	)

	redisAccountBroker := infra_event.NewRedisAccountNotificationBroker(
		redisClient,
		context.WithoutCancel(ctx),
		registry.Inject("Database").(infra_database.Database),
	)

	sendNotificationHandler := event_handlers.NewSendAccountWeatherNotification(
		infra_gateways.NewCPTECWeatherGateway(
			os.Getenv("WEATHER_API"),
			redisClient,
			context.WithoutCancel(ctx),
		),
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
