package registry_pattern

import (
	infra_repositories "weather_notification/src/modules/auth/infra/repositories"
	registry "weather_notification/src/modules/shared/infra"
	infra_database "weather_notification/src/modules/shared/infra/database"
)

func RegisterRepositories() {
	registry := registry.GetRegistryInstance()

	registry.Register(RepositoryKeys.Account, infra_repositories.NewMySQLAccountRepository(
		registry.Inject("Database").(infra_database.Database),
	))
	registry.Register(RepositoryKeys.Session, infra_repositories.NewMySQLSessionRepository(
		registry.Inject("Database").(infra_database.Database),
	))
}

type repositoryKeys struct {
	Account string
	Session string
}

var RepositoryKeys = repositoryKeys{Account: "AccountRepository", Session: "SessionRepository"}
