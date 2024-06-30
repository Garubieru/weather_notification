package registry_pattern

import (
	application_service "weather_notification/src/modules/auth/application/services"
	"weather_notification/src/modules/auth/domain/repositories"
	infra_services "weather_notification/src/modules/auth/infra/services"
	registry "weather_notification/src/modules/shared/infra"
)

func RegisterServices() {
	registry := registry.GetRegistryInstance()

	cryptoService := infra_services.NewBcryptCryptoService()

	registry.Register(ServiceKeys.CreateAccount, application_service.NewCreateAccountApplication(
		registry.Inject(RepositoryKeys.Account).(repositories.AccountRepository),
		cryptoService,
	))

	registry.Register(ServiceKeys.Login, application_service.NewLoginService(
		registry.Inject(RepositoryKeys.Account).(repositories.AccountRepository),
		registry.Inject(RepositoryKeys.Session).(repositories.SessionRepository),
		cryptoService,
	))

	registry.Register(ServiceKeys.Authenticate, application_service.NewAuthenticateSessionService(
		registry.Inject(RepositoryKeys.Session).(repositories.SessionRepository),
	))

	registry.Register(ServiceKeys.RetrieveAccount, application_service.NewRetrieveAccountInfo(
		registry.Inject(RepositoryKeys.Account).(repositories.AccountRepository),
	))
}

type serviceKeys struct {
	CreateAccount   string
	Login           string
	Authenticate    string
	RetrieveAccount string
}

var ServiceKeys = serviceKeys{
	CreateAccount:   "CreateAccountService",
	Login:           "LoginService",
	Authenticate:    "AuthenticateService",
	RetrieveAccount: "RetrieveAccount",
}
