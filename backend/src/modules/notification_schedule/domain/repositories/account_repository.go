package notification_schedule_repositories

import (
	notification_domain_entities "weather_notification/src/modules/notification_schedule/domain/entities"
	"weather_notification/src/modules/shared/value_objects"
)

type AccountRepository interface {
	FindById(id value_objects.ID) (*notification_domain_entities.AccountAggregateRoot, error)
	Save(account *notification_domain_entities.AccountAggregateRoot) error
}
