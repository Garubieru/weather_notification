package services_factories

import (
	application_service "weather_notification/src/modules/auth/application/services"
	"weather_notification/src/modules/auth/domain/repositories"
	registry_pattern "weather_notification/src/modules/auth/main/registry"
	registry "weather_notification/src/modules/shared/infra"
)

type ServiceFactories struct{}

func (f *ServiceFactories) CreateAuthenticateSession() application_service.AuthenticateSessionService {
	registry := registry.GetRegistryInstance()

	return application_service.NewAuthenticateSessionService(
		registry.Inject(registry_pattern.RepositoryKeys.Session).(repositories.SessionRepository),
	)
}
