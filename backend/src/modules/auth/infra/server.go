package infra

import (
	"fmt"
	"net/http"
	application_service "weather_notification/src/modules/auth/application/services"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Port               int
	Client             *gin.Engine
	authorizationGroup *gin.RouterGroup
}

func NewServer(port int) Server {
	client := gin.Default()
	authorizationGroup := client.Group("/")

	return Server{Port: port, Client: client, authorizationGroup: authorizationGroup}
}

func (server *Server) Listen() {
	server.Client.Run(fmt.Sprintf(":%d", server.Port))
}

func (server *Server) HealthCheck() {
	server.Client.Handle("GET", "/health-check", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"running": true})
	})
}

func (server *Server) IncludeRoute(command IncludeRouteCommand) {
	if !command.Private {
		server.Client.Handle(command.Method, command.Route, func(ctx *gin.Context) {
			command.Controller(ctx)
		})
		return
	}

	server.authorizationGroup.Handle(command.Method, command.Route, func(ctx *gin.Context) {
		command.Controller(ctx)
	})
}

func (server *Server) SetAuthentication(authenticationService application_service.AuthenticateSessionService) {
	server.authorizationGroup.Use(func(ctx *gin.Context) {
		sessionCookie, err := ctx.Cookie("session")

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized", "code": "Unauthorized"})
			ctx.Abort()
			return
		}

		output := authenticationService.Execute(application_service.AuthenticateSessionInput{SessionId: sessionCookie})

		if output.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not authenticate try again later", "code": "Authentication error"})
			ctx.Abort()
			return
		}

		if !output.Result.Authenticated {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized", "code": "Unauthorized"})
			ctx.Abort()
			return
		}

		ctx.Set("AccountId", output.Result.AccountId.Value)
		ctx.Next()
	})
}

type Controller func(ctx *gin.Context)

type IncludeRouteCommand struct {
	Method     string
	Route      string
	Controller Controller
	Private    bool
}
