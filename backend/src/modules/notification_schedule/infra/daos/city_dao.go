package infra_daos

import (
	"database/sql"
	"fmt"
	notification_schedule_daos "weather_notification/src/modules/notification_schedule/application/daos"
	infra_database "weather_notification/src/modules/shared/infra/database"
)

type CityMySQLDAO struct {
	database infra_database.Database
}

func (dao CityMySQLDAO) FindById(id string) (*notification_schedule_daos.CityDTO, error) {
	var cityDto notification_schedule_daos.CityDTO

	selectQuery := dao.database.QueryBuilder("city").SetColumns(
		[]string{"id", "external_id", "name", "state_code"},
	).Where("id = ?")

	row := dao.database.SelectOne(selectQuery.Select(), id)

	if err := row.Scan(&cityDto.Id, &cityDto.ExternalId, &cityDto.Name, &cityDto.StateCode); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query city `%s` error", id)
	}

	return &cityDto, nil
}

func NewCityMySQLDAO(database infra_database.Database) notification_schedule_daos.CityDAO {
	return CityMySQLDAO{database: database}
}
