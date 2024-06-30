package infra_daos

import (
	"database/sql"
	"time"
	notification_schedule_daos "weather_notification/src/modules/notification_schedule/application/daos"
	infra_database "weather_notification/src/modules/shared/infra/database"
)

type ScheduledNotificationMySQLDAO struct {
	database infra_database.Database
}

func (dao ScheduledNotificationMySQLDAO) FindByAccountId(accountId string) (
	[]notification_schedule_daos.ScheduledNotificationDTO,
	error,
) {
	selectQuery := dao.database.QueryBuilder("notification_schedule AS sn").
		SetColumns([]string{
			"sn.id",
			"sn.scheduled_date",
			"sn.interval_in_days",
			"sn.hour",
			"sn.is_coastal_city",
			"sn.active",
			"sn.method",
			"sn.account_id",
			"sn.city_id",
			"city.external_id AS city_external_id",
			"city.name AS city_name",
			"city.state_code AS city_state_code",
		}).
		Join("INNER JOIN city ON city.id = sn.city_id").
		Where("sn.account_id = ?").
		Select()

	rows, selectError := dao.database.Select(selectQuery, accountId)

	if selectError != nil {
		return nil, selectError
	}

	defer rows.Close()

	result := []notification_schedule_daos.ScheduledNotificationDTO{}

	for rows.Next() {
		var scheduledNotification ScheduledNotificationSchema

		if err := rows.Scan(
			&scheduledNotification.Id,
			&scheduledNotification.ScheduledDate,
			&scheduledNotification.IntervalInDays,
			&scheduledNotification.Hour,
			&scheduledNotification.IsCoastalCity,
			&scheduledNotification.Active,
			&scheduledNotification.Method,
			&scheduledNotification.AccountId,
			&scheduledNotification.CityId,
			&scheduledNotification.CityExternalId,
			&scheduledNotification.CityName,
			&scheduledNotification.CityStateCode,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}

		result = append(result, notification_schedule_daos.ScheduledNotificationDTO{
			Id:             scheduledNotification.Id,
			ScheduledDate:  scheduledNotification.ScheduledDate,
			Hour:           scheduledNotification.Hour,
			IntervalInDays: scheduledNotification.IntervalInDays,
			Method:         scheduledNotification.Method,
			IsCoastalCity:  scheduledNotification.IsCoastalCity,
			AccountId:      scheduledNotification.AccountId,
			Active:         scheduledNotification.Active,
			City: notification_schedule_daos.ScheduledNotificationCityDTO{
				Id:         scheduledNotification.CityId,
				ExternalId: scheduledNotification.CityExternalId,
				Name:       scheduledNotification.CityName,
				StateCode:  scheduledNotification.CityStateCode,
			},
		})
	}

	return result, nil
}

func (dao ScheduledNotificationMySQLDAO) FindByScheduledDate(time time.Time) (
	[]notification_schedule_daos.ScheduledNotificationDTO,
	error,
) {
	selectQuery := dao.database.QueryBuilder("notification_schedule AS sn").
		SetColumns([]string{
			"sn.id",
			"sn.scheduled_date",
			"sn.interval_in_days",
			"sn.hour",
			"sn.is_coastal_city",
			"sn.active",
			"sn.method",
			"sn.account_id",
			"sn.city_id",
			"city.external_id AS city_external_id",
			"city.name AS city_name",
			"city.state_code AS city_state_code",
		}).
		Join("INNER JOIN city ON city.id = sn.city_id").
		Where("sn.scheduled_date <= ? && sn.active = ?").
		Select()

	rows, selectError := dao.database.Select(selectQuery, time, true)

	if selectError != nil {
		return nil, selectError
	}

	defer rows.Close()

	result := []notification_schedule_daos.ScheduledNotificationDTO{}

	for rows.Next() {
		var scheduledNotification ScheduledNotificationSchema

		if err := rows.Scan(
			&scheduledNotification.Id,
			&scheduledNotification.ScheduledDate,
			&scheduledNotification.IntervalInDays,
			&scheduledNotification.Hour,
			&scheduledNotification.IsCoastalCity,
			&scheduledNotification.Active,
			&scheduledNotification.Method,
			&scheduledNotification.AccountId,
			&scheduledNotification.CityId,
			&scheduledNotification.CityExternalId,
			&scheduledNotification.CityName,
			&scheduledNotification.CityStateCode,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}

		result = append(result, notification_schedule_daos.ScheduledNotificationDTO{
			Id:             scheduledNotification.Id,
			ScheduledDate:  scheduledNotification.ScheduledDate,
			Hour:           scheduledNotification.Hour,
			IntervalInDays: scheduledNotification.IntervalInDays,
			Method:         scheduledNotification.Method,
			IsCoastalCity:  scheduledNotification.IsCoastalCity,
			Active:         scheduledNotification.Active,
			AccountId:      scheduledNotification.AccountId,
			City: notification_schedule_daos.ScheduledNotificationCityDTO{
				Id:         scheduledNotification.CityId,
				ExternalId: scheduledNotification.CityExternalId,
				Name:       scheduledNotification.CityName,
				StateCode:  scheduledNotification.CityStateCode,
			},
		})
	}

	return result, nil
}

func NewScheduledNotificationMySQLDAO(database infra_database.Database) notification_schedule_daos.ScheduledNotificationDAO {
	return ScheduledNotificationMySQLDAO{database: database}
}

type ScheduledNotificationSchema struct {
	Id             string
	ScheduledDate  time.Time
	IntervalInDays uint8
	Hour           uint8
	Method         string
	IsCoastalCity  bool
	Active         bool
	AccountId      string
	CityId         string
	CityExternalId string
	CityName       string
	CityStateCode  string
}
