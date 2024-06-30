package entities

import (
	"time"
	"weather_notification/src/modules/shared/value_objects"
)

type Session struct {
	Id         value_objects.ID
	AccountId  value_objects.ID
	ExpireDate time.Time
}

func (session Session) IsExpired() bool {
	return time.Now().After(session.ExpireDate)
}

func (session Session) GetId() string {
	return session.Id.Value
}

func NewSession(command SessionCreateCommand) Session {
	sessionTimeInMinutes := time.Duration(int8(30)) * time.Minute

	expireDate := time.Now()
	expireDate = expireDate.Add(sessionTimeInMinutes)

	return Session{AccountId: command.AccountId, Id: value_objects.NewID(), ExpireDate: expireDate}
}

func RecoverSession(command SessionRecoverCommand) *Session {
	return &Session{
		AccountId:  value_objects.ID{Value: command.AccountId},
		Id:         value_objects.ID{Value: command.Id},
		ExpireDate: command.ExpireDate,
	}
}

type SessionCreateCommand struct {
	AccountId value_objects.ID
}

type SessionRecoverCommand struct {
	Id         string
	AccountId  string
	ExpireDate time.Time
}
