package entities

import (
	"errors"
	"net/mail"
	crypto_service "weather_notification/src/modules/auth/domain/services"
	"weather_notification/src/modules/shared/value_objects"
)

type Account struct {
	Id              AccountId
	Name            string
	Username        string
	Email           Email
	Phone           string
	Password        string
	IsAuthenticated bool
}

type AccountId = value_objects.ID

func (account *Account) Authenticate(cryptoService crypto_service.CryptoService, password string) {
	account.IsAuthenticated = cryptoService.Verify(password, account.Password)
}

func (account Account) GetId() string {
	return account.Id.Value
}

func NewAccount(command CreateAccountCommand) (*Account, error) {
	email, err := NewEmail(command.Email)

	if err != nil {
		return nil, err
	}

	account := Account{
		Id:       value_objects.NewID(),
		Name:     command.Name,
		Username: command.Username,
		Email:    email,
		Phone:    command.Phone,
		Password: command.Password,
	}

	return &account, nil
}

func RecoverAccount(command RecoverAccountCommand) *Account {
	account := Account{
		Id:       value_objects.ID{Value: command.Id},
		Name:     command.Name,
		Username: command.Username,
		Email:    Email{Value: command.Email},
		Phone:    command.Phone,
		Password: command.Password,
	}

	return &account
}

type CreateAccountCommand struct {
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
}

type RecoverAccountCommand struct {
	Id       string
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
}

type Email struct {
	Value string
}

func NewEmail(email string) (Email, error) {
	result := Email{Value: email}

	isValid := result.validate()

	if !isValid {
		return Email{Value: ""}, errors.New("email")
	}

	return result, nil
}

func (e Email) validate() bool {
	_, err := mail.ParseAddress(e.Value)
	return err == nil
}
