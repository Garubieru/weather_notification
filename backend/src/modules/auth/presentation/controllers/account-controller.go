package controllers

import (
	"net/http"
	services "weather_notification/src/modules/auth/application/services"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	createAccountService services.CreateAccountService
	loginService         services.LoginService
	retrieveAccountInfo  services.RetrieveAccountInfo
}

func (controller AccountController) CreateAccount(ctx *gin.Context) {
	var accountRequest AccountPostRequest

	if err := ctx.ShouldBind(&accountRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output := controller.createAccountService.Execute(services.CreateAccountInput{
		Name:     accountRequest.Name,
		Username: accountRequest.Username,
		Email:    accountRequest.Email,
		Password: accountRequest.Password,
		Phone:    accountRequest.Phone})

	if output != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": output.Message, "code": output.Name})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

func (controller AccountController) Login(ctx *gin.Context) {
	var accountLoginRequest AccountLoginPostRequest

	if err := ctx.ShouldBind(&accountLoginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output := controller.loginService.Execute(services.LoginServiceInput{
		Username: accountLoginRequest.Username,
		Password: accountLoginRequest.Password,
	})

	if output.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": output.Error.Message, "code": output.Error.Name})
		return
	}

	ctx.SetCookie("session", output.Result.SessionId, 864000, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{"sessionId": output.Result.SessionId})
}

func (controller AccountController) GetSessionAccount(ctx *gin.Context) {
	accountId := ctx.GetString("AccountId")

	output := controller.retrieveAccountInfo.Execute(services.RetrieveAccountInfoInputDTO{
		AccountId: accountId,
	})

	if output.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": output.Error.Message, "code": output.Error.Name})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":       output.Result.Id,
		"name":     output.Result.Name,
		"username": output.Result.Username,
		"email":    output.Result.Email,
		"phone":    output.Result.Phone,
	})
}

type AccountPostRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type AccountLoginPostRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewAccountController(input NewAccountControllerInput) AccountController {
	return AccountController{
		createAccountService: input.CreateAccountService,
		loginService:         input.LoginService,
		retrieveAccountInfo:  input.RetrieveAccountInfo,
	}
}

type NewAccountControllerInput struct {
	CreateAccountService services.CreateAccountService
	LoginService         services.LoginService
	RetrieveAccountInfo  services.RetrieveAccountInfo
}
