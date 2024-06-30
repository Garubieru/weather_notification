package infra_repositories

import (
	"database/sql"
	"fmt"
	"weather_notification/src/modules/auth/domain/entities"
	infra_database "weather_notification/src/modules/shared/infra/database"
	"weather_notification/src/modules/shared/value_objects"
)

type MySQLAccountRepository struct {
	database infra_database.Database
}

func (repository MySQLAccountRepository) Save(account *entities.Account) error {
	if account == nil {
		return fmt.Errorf("accountRepository.Save: account must not be empty")
	}

	err := repository.database.Exec(`INSERT INTO account (id, name, username, password, email, phone) VALUES (?, ?, ?, ?, ?, ?)`,
		account.Id.Value,
		account.Name,
		account.Username,
		account.Password,
		account.Email.Value,
		account.Phone,
	)

	if err != nil {
		return fmt.Errorf("accountRepository.Save %v", err)
	}

	return nil
}

func (repository MySQLAccountRepository) FindById(id value_objects.ID) (*entities.Account, error) {
	row := repository.database.SelectOne("SELECT id, name, username, password, email, phone FROM account WHERE id = ?", id.Value)

	var account AccountSchema

	if err := row.Scan(&account.Id, &account.Name, &account.Username, &account.Password, &account.Email, &account.Phone); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("accountRepository.FindByUsername %s, error: %v", id.Value, err)
	}

	return entities.RecoverAccount(entities.RecoverAccountCommand{
		Id:       account.Id,
		Name:     account.Name,
		Username: account.Username,
		Email:    account.Email,
		Phone:    account.Phone,
		Password: account.Password,
	}), nil
}

func (repository MySQLAccountRepository) FindByUsername(username string) (*entities.Account, error) {
	row := repository.database.SelectOne("SELECT id, name, username, password, email, phone FROM account WHERE username = ?", username)

	var account AccountSchema

	if err := row.Scan(&account.Id, &account.Name, &account.Username, &account.Password, &account.Email, &account.Phone); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("accountRepository.FindByUsername %s, error: %v", username, err)
	}

	return entities.RecoverAccount(entities.RecoverAccountCommand{
		Id:       account.Id,
		Name:     account.Name,
		Username: account.Username,
		Email:    account.Email,
		Phone:    account.Phone,
		Password: account.Password,
	}), nil
}

func (repository MySQLAccountRepository) FindByEmail(email string) (*entities.Account, error) {
	row := repository.database.SelectOne("SELECT id, name, username, password, email, phone FROM account WHERE email = ?", email)

	var account AccountSchema

	if err := row.Scan(&account.Id, &account.Name, &account.Username, &account.Password, &account.Email, &account.Phone); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("accountRepository.FindByEmail %s, error: %v", email, err)
	}

	return entities.RecoverAccount(entities.RecoverAccountCommand{
		Id:       account.Id,
		Name:     account.Name,
		Username: account.Username,
		Email:    account.Email,
		Phone:    account.Phone,
		Password: account.Password,
	}), nil
}

type AccountSchema struct {
	Id       string
	Name     string
	Username string
	Password string
	Email    string
	Phone    string
}

func NewMySQLAccountRepository(database infra_database.Database) MySQLAccountRepository {
	return MySQLAccountRepository{database: database}
}
