package notification_registry_pattern

import (
	infra_daos "weather_notification/src/modules/notification_schedule/infra/daos"
	registry "weather_notification/src/modules/shared/infra"
	infra_database "weather_notification/src/modules/shared/infra/database"
)

func RegisterNotificationScheduleDAOs() {
	registry := registry.GetRegistryInstance()

	registry.Register(DAOKeys.ScheduledNotificationDAO, infra_daos.NewScheduledNotificationMySQLDAO(
		registry.Inject("Database").(infra_database.Database),
	))

	registry.Register(DAOKeys.CityDAO, infra_daos.NewCityMySQLDAO(
		registry.Inject("Database").(infra_database.Database),
	))

	registry.Register(DAOKeys.NotificationDAO, infra_daos.NewNotificationMySQLDAO(
		registry.Inject("Database").(infra_database.Database),
	))
}

type daosKeys struct {
	ScheduledNotificationDAO string
	CityDAO                  string
	NotificationDAO          string
}

var DAOKeys = daosKeys{
	ScheduledNotificationDAO: "ScheduledNotificationDAO",
	CityDAO:                  "CityDAO",
	NotificationDAO:          "NotificationDAO",
}
