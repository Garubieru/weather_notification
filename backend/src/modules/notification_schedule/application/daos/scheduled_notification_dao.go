package notification_schedule_daos

import "time"

type ScheduledNotificationDAO interface {
	FindByAccountId(accountId string) ([]ScheduledNotificationDTO, error)
	FindByScheduledDate(date time.Time) ([]ScheduledNotificationDTO, error)
}

type ScheduledNotificationDTO struct {
	Id             string                       `json:"id"`
	ScheduledDate  time.Time                    `json:"scheduledDate"`
	IntervalInDays uint8                        `json:"intervalInDays"`
	Hour           uint8                        `json:"hour"`
	Method         string                       `json:"method"`
	IsCoastalCity  bool                         `json:"isCoastalCity"`
	Active         bool                         `json:"active"`
	City           ScheduledNotificationCityDTO `json:"city"`
	AccountId      string                       `json:"accountId"`
}

type ScheduledNotificationCityDTO struct {
	Id         string `json:"id"`
	ExternalId string `json:"externalId"`
	Name       string `json:"name"`
	StateCode  string `json:"stateCode"`
}
