package infra_jobs

import notification_schedule_services "weather_notification/src/modules/notification_schedule/application/services"

type SendScheduledNotificationJob struct {
	service notification_schedule_services.SendScheduledNotification
}

func (job SendScheduledNotificationJob) Handle() error {
	return job.service.Execute()
}

func NewSendScheduleNotification(service notification_schedule_services.SendScheduledNotification) Job {
	return SendScheduledNotificationJob{service: service}
}

