package notification_registry_pattern

import (
	notification_schedule_daos "weather_notification/src/modules/notification_schedule/application/daos"
	notification_schedule_query "weather_notification/src/modules/notification_schedule/application/queries"
	notification_schedule_services "weather_notification/src/modules/notification_schedule/application/services"
	event_broker "weather_notification/src/modules/notification_schedule/domain/event"
	notification_schedule_repositories "weather_notification/src/modules/notification_schedule/domain/repositories"
	registry "weather_notification/src/modules/shared/infra"
)

func RegisterNotificationScheduleServices() {
	registry := registry.GetRegistryInstance()

	registry.Register(NotificationServiceKeys.ScheduleWeatherNotification,
		notification_schedule_services.NewScheduleWeatherNotification(
			registry.Inject(RepositoryKeys.AccountNotificationScheduleRepository).(notification_schedule_repositories.AccountRepository),
			registry.Inject(DAOKeys.CityDAO).(notification_schedule_daos.CityDAO),
		))

	registry.Register(NotificationServiceKeys.DeactivateWeatherNotificationSchedule,
		notification_schedule_services.NewDeactivateWeatherNotificationScheduleService(
			registry.Inject(RepositoryKeys.AccountNotificationScheduleRepository).(notification_schedule_repositories.AccountRepository),
		))

	registry.Register(NotificationServiceKeys.ActivateWeatherNotificationSchedule,
		notification_schedule_services.NewActivateWeatherNotificationScheduleService(
			registry.Inject(RepositoryKeys.AccountNotificationScheduleRepository).(notification_schedule_repositories.AccountRepository),
		))

	registry.Register(NotificationServiceKeys.ListAccountWeatherNotificationSchedules,
		notification_schedule_query.NewListAccountWeatherNotificationService(
			registry.Inject(DAOKeys.ScheduledNotificationDAO).(notification_schedule_daos.ScheduledNotificationDAO),
		),
	)

	registry.Register(NotificationServiceKeys.SendScheduledNotifications,
		notification_schedule_services.NewSendScheduleNotification(
			registry.Inject(DAOKeys.ScheduledNotificationDAO).(notification_schedule_daos.ScheduledNotificationDAO),
			registry.Inject(InfraKeys.KafkaEventBroker).(event_broker.EventBroker),
		),
	)
}

type notificationServiceKeys struct {
	ScheduleWeatherNotification             string
	DeactivateWeatherNotificationSchedule   string
	ActivateWeatherNotificationSchedule     string
	ListAccountWeatherNotificationSchedules string
	SendScheduledNotifications              string
}

var NotificationServiceKeys = notificationServiceKeys{
	ScheduleWeatherNotification:             "ScheduleWeatherNotification",
	DeactivateWeatherNotificationSchedule:   "DeactivateWeatherNotificationSchedule",
	ListAccountWeatherNotificationSchedules: "ListAccountWeatherNotificationSchedules",
	SendScheduledNotifications:              "SendScheduledNotifications",
}
