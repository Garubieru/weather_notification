package mocks

import (
	entities "weather_notification/src/modules/auth/domain/entities"
	"weather_notification/src/modules/shared/value_objects"
)

type AccountRepositoryInMemory struct {
	Accounts map[string]entities.Account
}

func (r *AccountRepositoryInMemory) FindById(id value_objects.ID) (*entities.Account, error) {
	if account, ok := r.Accounts[id.Value]; ok {
		return &account, nil
	}

	return nil, nil
}

func (r *AccountRepositoryInMemory) FindByEmail(email string) (*entities.Account, error) {
	for _, account := range r.Accounts {
		if account.Email.Value == email {
			return &account, nil
		}
	}

	return nil, nil
}

func (r *AccountRepositoryInMemory) FindByUsername(username string) (*entities.Account, error) {
	for _, account := range r.Accounts {
		if account.Username == username {
			return &account, nil
		}
	}

	return nil, nil
}

func (r *AccountRepositoryInMemory) Save(account *entities.Account) error {
	r.Accounts[account.Id.Value] = *account
	return nil
}

func NewAccountRepositoryInMemory() *AccountRepositoryInMemory {
	return &AccountRepositoryInMemory{Accounts: make(map[string]entities.Account)}
}
