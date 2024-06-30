package auth_application_service

import (
	"errors"
	"fmt"
	entities "weather_notification/src/modules/auth/domain/entities"
	repositories "weather_notification/src/modules/auth/domain/repositories"
	services "weather_notification/src/modules/auth/domain/services"
	application_service "weather_notification/src/modules/shared/application"
	application_error "weather_notification/src/modules/shared/application/errors"
)

type CreateAccountApplicationImp struct {
	accountRepository repositories.AccountRepository
	cryptoService     services.CryptoService
}

func (service CreateAccountApplicationImp) Execute(input CreateAccountInput) *application_error.ApplicationError {
	err := service.validateAccount(input.Username, input.Email)

	if err != nil {
		return err
	}

	hash, encryptError := service.cryptoService.Encrypt(input.Password)

	if encryptError != nil {
		return &application_error.ApplicationError{Message: "Error on password encryption", Err: encryptError, Name: "PasswordEncryptionError"}
	}

	account, validationError := entities.NewAccount(entities.CreateAccountCommand{
		Name:     input.Name,
		Username: input.Username,
		Email:    input.Email,
		Phone:    input.Phone,
		Password: hash,
	})

	if validationError != nil {
		return &application_error.ApplicationError{Message: fmt.Sprintf("Invalid %s", validationError.Error()), Err: err, Name: "ValidationError"}
	}

	if err := service.accountRepository.Save(account); err != nil {
		return &application_error.ApplicationError{
			Message: "Failed to save account",
			Name:    "AccountRepositoryError",
			Err:     err}
	}

	return nil
}

func (r CreateAccountApplicationImp) validateAccount(username string, email string) *application_error.ApplicationError {
	foundAccount, err := r.accountRepository.FindByUsername(username)

	if err != nil {
		return application_error.NewRepositoryError("AccountRepository")
	}

	if foundAccount != nil {
		return &application_error.ApplicationError{
			Message: "Account username already taken",
			Name:    "UsernameAlreadyTaken",
			Err:     errors.New("account username already taken")}
	}

	foundByEmailAccount, err := r.accountRepository.FindByEmail(email)

	if err != nil {
		return application_error.NewRepositoryError("AccountRepository")
	}

	if foundByEmailAccount != nil {
		return &application_error.ApplicationError{
			Message: "Account email already exists",
			Name:    "EmailAlreadyTaken",
			Err:     errors.New("account email already taken")}
	}

	return nil
}

func NewCreateAccountApplication(accountRepository repositories.AccountRepository, cryptoSevice services.CryptoService) CreateAccountService {
	return &CreateAccountApplicationImp{accountRepository: accountRepository, cryptoService: cryptoSevice}
}

type CreateAccountInput struct {
	Name     string
	Username string
	Email    string
	Password string
	Phone    string
}

type CreateAccountService = application_service.ApplicationService[CreateAccountInput, *application_error.ApplicationError]
