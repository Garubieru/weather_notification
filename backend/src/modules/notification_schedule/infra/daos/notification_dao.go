package infra_daos

import (
	"database/sql"
	"encoding/json"
	"errors"
	notification_schedule_daos "weather_notification/src/modules/notification_schedule/application/daos"
	event_broker "weather_notification/src/modules/notification_schedule/domain/event"
	infra_database "weather_notification/src/modules/shared/infra/database"
)

type NotificationMySQLDAO struct {
	database infra_database.Database
}

func NewNotificationMySQLDAO(database infra_database.Database) notification_schedule_daos.NotificationDAO {
	return NotificationMySQLDAO{database: database}
}

func (dao NotificationMySQLDAO) FindNotifications(accountId string) ([]notification_schedule_daos.NotificationDTO, error) {
	countQueryBuilder := dao.database.QueryBuilder("notification").SetColumns(
		[]string{"COUNT(id)"},
	).Where("account_id = ?")

	var length int

	if err := dao.database.SelectOne(countQueryBuilder.Select(), accountId).Scan(&length); err != nil {
		if err == sql.ErrNoRows {
			return []notification_schedule_daos.NotificationDTO{}, nil
		}
		return nil, err
	}

	queryBuilder := dao.database.QueryBuilder("notification").
		SetColumns(
			[]string{"payload"},
		).
		Where("account_id = ?").
		OrderBy("created_at DESC")

	rows, err := dao.database.Select(queryBuilder.Select(), accountId)

	if err != nil {
		return nil, err
	}

	result := make([]event_broker.Event, 0, length)

	defer rows.Close()

	for rows.Next() {
		var notificationDTOValue []byte

		if err := rows.Scan(&notificationDTOValue); err != nil {
			if err == sql.ErrNoRows {
				return result, nil
			}
			return nil, err
		}

		var notificationDTO notification_schedule_daos.NotificationDTO

		if err := json.Unmarshal(notificationDTOValue, &notificationDTO); err != nil {
			return nil, errors.New("could not parse data")
		}

		result = append(result, notificationDTO)
	}

	return result, nil
}
