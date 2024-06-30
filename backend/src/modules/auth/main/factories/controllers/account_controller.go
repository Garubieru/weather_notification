package controllers_factories

import (
	auth_application_service "weather_notification/src/modules/auth/application/services"
	registry_pattern "weather_notification/src/modules/auth/main/registry"
	"weather_notification/src/modules/auth/presentation/controllers"
	registry "weather_notification/src/modules/shared/infra"
)

func NewAccountControllerFactory() controllers.AccountController {
	registry := registry.GetRegistryInstance()

	accountController := controllers.NewAccountController(controllers.NewAccountControllerInput{
		CreateAccountService: registry.Inject(registry_pattern.ServiceKeys.CreateAccount).(auth_application_service.CreateAccountService),
		LoginService:         registry.Inject(registry_pattern.ServiceKeys.Login).(auth_application_service.LoginService),
		RetrieveAccountInfo:  registry.Inject(registry_pattern.ServiceKeys.RetrieveAccount).(auth_application_service.RetrieveAccountInfo),
	})

	return accountController
}
