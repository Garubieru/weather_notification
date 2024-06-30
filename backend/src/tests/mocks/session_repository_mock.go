package mocks

import (
	entities "weather_notification/src/modules/auth/domain/entities"
	"weather_notification/src/modules/shared/value_objects"
)

type SessionRepositoryMemory struct {
	Sessions map[string]*entities.Session
}

func (r *SessionRepositoryMemory) FindById(id value_objects.ID) (*entities.Session, error) {
	session, ok := r.Sessions[id.Value]

	if ok {
		return session, nil
	}

	return nil, nil
}

func (r *SessionRepositoryMemory) Save(session *entities.Session) error {
	r.Sessions[session.Id.Value] = session
	return nil
}

func NewSessionRepositoryMemory() *SessionRepositoryMemory {
	return &SessionRepositoryMemory{Sessions: make(map[string]*entities.Session)}
}
