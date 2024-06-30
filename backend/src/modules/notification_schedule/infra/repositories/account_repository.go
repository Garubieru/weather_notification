package notification_schedule_repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
	notification_domain_entities "weather_notification/src/modules/notification_schedule/domain/entities"
	infra_database "weather_notification/src/modules/shared/infra/database"
	"weather_notification/src/modules/shared/value_objects"
)

type MySQLAccountRepository struct {
	database infra_database.Database
}

func (repository MySQLAccountRepository) FindById(id value_objects.ID) (*notification_domain_entities.AccountAggregateRoot, error) {
	var count int

	repository.database.SelectOne(`
		SELECT COUNT(id)
		FROM notification_schedule
		WHERE account_id = ? && active = TRUE`, id.Value).Scan(&count)

	rows, err := repository.database.Select(`
		SELECT
			account.id,
			account.email,
			account.phone,
			sn.id as schedule_id,
			sn.scheduled_date,
			sn.interval_in_days,
			sn.hour,
			sn.city_id,
			sn.is_coastal_city,
			sn.active,
			sn.method
		FROM
			account
		LEFT JOIN notification_schedule AS sn ON sn.account_id = account.id
		WHERE
			account.id = ?
	`, id.Value)

	if err != nil {
		return nil, fmt.Errorf("account query failed")
	}

	scheduledNotifications := make([]AccountScheduledNotificationSchema, 0, count)

	defer rows.Close()

	for rows.Next() {
		var scheduledNotification AccountScheduledNotificationSchema

		if err := rows.Scan(
			&scheduledNotification.Id,
			&scheduledNotification.Email,
			&scheduledNotification.Phone,
			&scheduledNotification.ScheduleId,
			&scheduledNotification.ScheduledDate,
			&scheduledNotification.IntervalInDays,
			&scheduledNotification.Hour,
			&scheduledNotification.CityId,
			&scheduledNotification.IsCoastalCity,
			&scheduledNotification.Active,
			&scheduledNotification.Method,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, fmt.Errorf("error when trying to recover values %v", err)
		}

		scheduledNotifications = append(scheduledNotifications, scheduledNotification)
	}

	accountInfo := scheduledNotifications[0]

	scheduleNotificationsCommand := make(
		[]notification_domain_entities.RecoverWeatherNotificationScheduleCommand,
		0,
		len(scheduledNotifications))

	for _, scheduledNotification := range scheduledNotifications {
		if scheduledNotification.ScheduleId == nil {
			continue
		}

		scheduleNotificationCommand := notification_domain_entities.RecoverWeatherNotificationScheduleCommand{
			Id:             *scheduledNotification.ScheduleId,
			ScheduledDate:  *scheduledNotification.ScheduledDate,
			CityId:         *scheduledNotification.CityId,
			IsCityCoastal:  *scheduledNotification.IsCoastalCity,
			Active:         *scheduledNotification.Active,
			IntervalInDays: int8(*scheduledNotification.IntervalInDays),
			Hour:           int8(*scheduledNotification.Hour),
			Method:         *scheduledNotification.Method,
		}

		scheduleNotificationsCommand = append(scheduleNotificationsCommand, scheduleNotificationCommand)
	}

	output := notification_domain_entities.RecoverAccountAggregate(
		notification_domain_entities.RecoverAccountAggregateCommand{
			Id:                     accountInfo.Id,
			Email:                  accountInfo.Email,
			Phone:                  accountInfo.Phone,
			ScheduledNotifications: scheduleNotificationsCommand,
		},
	)

	return &output, nil
}

func (repository MySQLAccountRepository) Save(account *notification_domain_entities.AccountAggregateRoot) error {
	if account == nil {
		return nil
	}

	repository.database.Transaction(func(tx *sql.Tx) error {
		baseQuery := repository.database.
			QueryBuilder("notification_schedule").
			SetColumns([]string{"id", "scheduled_date", "interval_in_days", "hour", "city_id", "is_coastal_city", "active", "method", "account_id"})

		if len(account.ScheduledNotifications.GetNewItems()) > 0 {
			insertValues := make([]interface{}, 0, len(account.ScheduledNotifications.GetNewItems()))

			for _, notification := range account.ScheduledNotifications.GetNewItems() {
				insertValues = append(insertValues,
					notification.Id.Value,
					notification.ScheduledDate,
					notification.IntervalInDays,
					notification.Hour.Value,
					notification.City.Id,
					notification.City.IsCoastal,
					notification.Active,
					notification.Method,
					account.Id.Value,
				)
			}

			if _, insertError := tx.Exec(baseQuery.Insert(len(account.ScheduledNotifications.GetNewItems())), insertValues...); insertError != nil {
				return insertError
			}
		}

		if len(account.ScheduledNotifications.GetDirtyItems()) > 0 {
			baseQuery.SetColumns([]string{"scheduled_date", "interval_in_days", "hour", "city_id", "is_coastal_city", "active", "method"})
			baseQuery.Where("id = ?")
			updateQuery := baseQuery.Update()

			for _, notification := range account.ScheduledNotifications.GetDirtyItems() {
				if _, updatedErr := tx.Exec(updateQuery,
					notification.ScheduledDate,
					notification.IntervalInDays,
					notification.Hour.Value,
					notification.City.Id,
					notification.City.IsCoastal,
					notification.Active,
					notification.Method,
					notification.Id.Value,
				); updatedErr != nil {
					return updatedErr
				}
			}
		}

		if len(account.ScheduledNotifications.GetRemovedItems()) > 0 {
			toDelete := make([]string, 0, len(account.ScheduledNotifications.GetRemovedItems()))
			placeholders := make([]string, 0, len(account.ScheduledNotifications.GetRemovedItems()))

			for _, notification := range account.ScheduledNotifications.GetRemovedItems() {
				toDelete = append(toDelete, notification.Id.Value)
				placeholders = append(placeholders, "?")
			}

			baseQuery.Where(fmt.Sprintf("id IN (%s)", strings.Join(placeholders, ", ")))

			deleteQuery := baseQuery.Delete()

			if _, err := tx.Exec(deleteQuery, toDelete); err != nil {
				return err
			}
		}

		return nil
	})

	return nil
}

type AccountScheduledNotificationSchema struct {
	Id             string
	Email          string
	Phone          string
	ScheduleId     *string
	ScheduledDate  *time.Time
	IntervalInDays *uint8
	Hour           *uint8
	CityId         *string
	IsCoastalCity  *bool
	Active         *bool
	Method         *string
}

func NewAccountScheduleMySQLRepository(database infra_database.Database) MySQLAccountRepository {
	return MySQLAccountRepository{database: database}
}
