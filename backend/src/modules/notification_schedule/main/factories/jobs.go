package notification_schedule_factory

import (
	notification_schedule_services "weather_notification/src/modules/notification_schedule/application/services"
	infra_jobs "weather_notification/src/modules/notification_schedule/infra/jobs"
	notification_registry_pattern "weather_notification/src/modules/notification_schedule/main/registry"
	registry "weather_notification/src/modules/shared/infra"
)

func StartJobs() {
	registry := registry.GetRegistryInstance()
	jobScheduler := infra_jobs.JobScheduler{}

	sendScheduledNotifcationsJob := infra_jobs.NewSendScheduleNotification(
		registry.Inject(notification_registry_pattern.NotificationServiceKeys.SendScheduledNotifications).(notification_schedule_services.SendScheduledNotification),
	)

	jobScheduler.Schedule("send_notifications", "*/10 * * * * *", sendScheduledNotifcationsJob)
}
