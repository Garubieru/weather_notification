package notification_schedule_factory

import (
	"context"
	notification_schedule_query "weather_notification/src/modules/notification_schedule/application/queries"
	notification_schedule_services "weather_notification/src/modules/notification_schedule/application/services"
	event_broker "weather_notification/src/modules/notification_schedule/domain/event"
	notification_registry_pattern "weather_notification/src/modules/notification_schedule/main/registry"
	schedule_controller "weather_notification/src/modules/notification_schedule/presentation/controllers"
	registry "weather_notification/src/modules/shared/infra"

	"github.com/go-redis/redis/v8"
)

func NewScheduleController() schedule_controller.ScheduleController {
	registry := registry.GetRegistryInstance()

	return *schedule_controller.NewScheduleController(
		registry.Inject(notification_registry_pattern.NotificationServiceKeys.ScheduleWeatherNotification).(notification_schedule_services.ScheduleNotificationApplicationService),
		registry.Inject(notification_registry_pattern.NotificationServiceKeys.DeactivateWeatherNotificationSchedule).(notification_schedule_services.DeactivateWeatherNotificationScheduleService),
		registry.Inject(notification_registry_pattern.NotificationServiceKeys.ActivateWeatherNotificationSchedule).(notification_schedule_services.ActivateWeatherNotificationScheduleService),
		registry.Inject(notification_registry_pattern.NotificationServiceKeys.ListAccountWeatherNotificationSchedules).(notification_schedule_query.ListAccountWeatherNotificationsService),
	)
}

func NewStreamController(ctx context.Context) schedule_controller.StreamController {
	registry := registry.GetRegistryInstance()

	return schedule_controller.NewStreamController(
		registry.Inject(notification_registry_pattern.InfraKeys.RedisClient).(redis.Client),
		registry.Inject(notification_registry_pattern.InfraKeys.RedisEventBroker).(event_broker.EventBroker),
	)
}
