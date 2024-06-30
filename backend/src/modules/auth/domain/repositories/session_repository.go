package repositories

import (
	entities "weather_notification/src/modules/auth/domain/entities"
	"weather_notification/src/modules/shared/value_objects"
)

type SessionRepository interface {
	Save(session *entities.Session) error
	FindById(sessionId value_objects.ID) (*entities.Session, error)
	Delete(session *entities.Session) error
}
