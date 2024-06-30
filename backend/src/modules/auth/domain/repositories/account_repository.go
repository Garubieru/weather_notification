package repositories

import (
	account "weather_notification/src/modules/auth/domain/entities"
	"weather_notification/src/modules/shared/value_objects"
)

type AccountRepository interface {
	FindByUsername(username string) (*account.Account, error)
	FindByEmail(email string) (*account.Account, error)
	FindById(id value_objects.ID) (*account.Account, error)
	Save(account *account.Account) error
}
