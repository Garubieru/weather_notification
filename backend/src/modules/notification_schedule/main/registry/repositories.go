package notification_registry_pattern

import (
	notification_schedule_repository "weather_notification/src/modules/notification_schedule/infra/repositories"
	registry "weather_notification/src/modules/shared/infra"
	infra_database "weather_notification/src/modules/shared/infra/database"
)

func RegisterNotificationRepositories() {
	registry := registry.GetRegistryInstance()

	registry.Register(RepositoryKeys.AccountNotificationScheduleRepository, notification_schedule_repository.NewAccountScheduleMySQLRepository(
		registry.Inject("Database").(infra_database.Database),
	))

}

type repositoryKeys struct {
	AccountNotificationScheduleRepository string
}

var RepositoryKeys = repositoryKeys{AccountNotificationScheduleRepository: "AccountNotificationScheduleRepository"}
