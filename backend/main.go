package main

import (
	"context"
	"log"
	"os"
	"strconv"
	application_service "weather_notification/src/modules/auth/application/services"
	"weather_notification/src/modules/auth/infra"
	controllers_factories "weather_notification/src/modules/auth/main/factories/controllers"
	registry_pattern "weather_notification/src/modules/auth/main/registry"
	notification_schedule_factory "weather_notification/src/modules/notification_schedule/main/factories"
	notification_registry_pattern "weather_notification/src/modules/notification_schedule/main/registry"
	registry "weather_notification/src/modules/shared/infra"
	infra_database "weather_notification/src/modules/shared/infra/database"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	serverPort, err := strconv.ParseUint(os.Getenv("SERVER_PORT"), 10, 16)

	if err != nil {
		log.Fatal("Provide a valid server port number")
	}

	server := infra.NewServer(int(serverPort))

	context := context.Background()

	registry := registry.GetRegistryInstance()

	port, err := strconv.ParseUint(os.Getenv("DATABASE_PORT"), 10, 16)

	if err != nil {
		log.Fatal("Provide a valid port number")
	}

	database := infra_database.NewMySQLDatabase(infra_database.MySQLDatabaseConfig{
		User:     os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		Name:     os.Getenv("DATABASE_NAME"),
		Host:     os.Getenv("DATABASE_HOST"),
		Port:     int(port),
	})

	database.Connect()

	registry.Register("Database", database)

	registry_pattern.RegisterRepositories()
	registry_pattern.RegisterServices()

	notification_registry_pattern.RegisterInfra(context)
	notification_registry_pattern.RegisterNotificationScheduleDAOs()
	notification_registry_pattern.RegisterNotificationRepositories()
	notification_registry_pattern.RegisterNotificationScheduleServices()

	server.SetAuthentication(registry.Inject(registry_pattern.ServiceKeys.Authenticate).(application_service.AuthenticateSessionService))

	accountController := controllers_factories.NewAccountControllerFactory()
	scheduleController := notification_schedule_factory.NewScheduleController()
	streamController := notification_schedule_factory.NewStreamController(context)

	notification_schedule_factory.StartJobs()

	server.IncludeRoute(infra.IncludeRouteCommand{
		Method:     "POST",
		Route:      "/v1/register",
		Controller: accountController.CreateAccount,
		Private:    false,
	})

	server.IncludeRoute(infra.IncludeRouteCommand{
		Method:     "POST",
		Route:      "/v1/login",
		Controller: accountController.Login,
		Private:    false,
	})

	server.IncludeRoute(infra.IncludeRouteCommand{
		Method:     "GET",
		Route:      "/v1/session",
		Controller: accountController.GetSessionAccount,
		Private:    true,
	})

	server.IncludeRoute(infra.IncludeRouteCommand{
		Method:     "POST",
		Route:      "/v1/account/schedules",
		Controller: scheduleController.Schedule,
		Private:    true,
	})

	server.IncludeRoute(infra.IncludeRouteCommand{
		Method:     "DELETE",
		Route:      "/v1/account/schedules/:scheduleId",
		Controller: scheduleController.DeactivateSchedule,
		Private:    true,
	})

	server.IncludeRoute(infra.IncludeRouteCommand{
		Method:     "PATCH",
		Route:      "/v1/account/schedules/:scheduleId",
		Controller: scheduleController.ActivateSchedule,
		Private:    true,
	})

	server.IncludeRoute(infra.IncludeRouteCommand{
		Method:     "GET",
		Route:      "/v1/account/schedules",
		Controller: scheduleController.ListAccountSchedules,
		Private:    true,
	})

	server.IncludeRoute(infra.IncludeRouteCommand{
		Method:     "GET",
		Route:      "/v1/account/notifications",
		Controller: scheduleController.ListAccountNotifications,
		Private:    true,
	})

	server.IncludeRoute(infra.IncludeRouteCommand{
		Method:     "GET",
		Route:      "/v1/stream",
		Controller: streamController.StartStream,
		Private:    true,
	})

	server.HealthCheck()

	server.Listen()
}
