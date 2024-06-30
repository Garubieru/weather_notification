package infra_repositories

import (
	"database/sql"
	"fmt"
	"time"
	"weather_notification/src/modules/auth/domain/entities"
	infra_database "weather_notification/src/modules/shared/infra/database"
	"weather_notification/src/modules/shared/value_objects"
)

type MySQLSessionRepository struct {
	database infra_database.Database
}

func (repository MySQLSessionRepository) Save(session *entities.Session) error {
	if session == nil {
		return fmt.Errorf("sessionRepository.Save: session must not be empty")
	}

	err := repository.database.Exec(`INSERT INTO session (id, expire_date, account_id) VALUES (?, ?, ?)`,
		session.Id.Value,
		session.ExpireDate,
		session.AccountId.Value,
	)

	if err != nil {
		return fmt.Errorf("sessionRepository.Save %v", err)
	}

	return nil
}

func (repository MySQLSessionRepository) Delete(session *entities.Session) error {
	if session == nil {
		return nil
	}

	err := repository.database.Exec(`DELETE FROM session WHERE id = ?`, session.Id.Value)

	if err != nil {
		return fmt.Errorf("sessionRepository.Delete %v", err)
	}

	return nil
}

func (repository MySQLSessionRepository) FindById(id value_objects.ID) (*entities.Session, error) {
	row := repository.database.SelectOne("SELECT id, account_id, expire_date FROM session WHERE id = ?", id.Value)

	var sessionSchema SessionSchema

	if err := row.Scan(&sessionSchema.Id, &sessionSchema.AccountId, &sessionSchema.ExpireDate); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("sessionRepository.FindById %s, error: %v", id.Value, err)
	}

	return entities.RecoverSession(entities.SessionRecoverCommand{
		Id:         sessionSchema.Id,
		AccountId:  sessionSchema.AccountId,
		ExpireDate: sessionSchema.ExpireDate,
	}), nil
}

type SessionSchema struct {
	Id         string
	AccountId  string
	ExpireDate time.Time
}

func NewMySQLSessionRepository(database infra_database.Database) MySQLSessionRepository {
	return MySQLSessionRepository{database: database}
}
