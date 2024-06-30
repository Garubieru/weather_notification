package schedule_controllers

import (
	"net/http"
	"time"
	notification_schedule_query "weather_notification/src/modules/notification_schedule/application/queries"
	notification_schedule_services "weather_notification/src/modules/notification_schedule/application/services"

	"github.com/gin-gonic/gin"
)

type ScheduleController struct {
	scheduleWeatherNotification             notification_schedule_services.ScheduleNotificationApplicationService
	deactivateWeatherNotification           notification_schedule_services.DeactivateWeatherNotificationScheduleService
	activateWeatherNotification             notification_schedule_services.ActivateWeatherNotificationScheduleService
	listAccountWeatherNotificationSchedules notification_schedule_query.ListAccountWeatherNotificationsService
}

func (controller *ScheduleController) Schedule(ctx *gin.Context) {
	var request ScheduleControllerRequest

	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accountId := ctx.GetString("AccountId")

	output := controller.scheduleWeatherNotification.Execute(notification_schedule_services.ScheduleNotificationInputDTO{
		AccountId:      accountId,
		CityId:         request.CityId,
		IsCityCoastal:  request.IsCoastalCity,
		IntervalInDays: request.IntervalInDays,
		Hour:           request.Hour,
		Method:         request.Method,
	})

	if output.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": output.Error.Message, "code": output.Error.Name})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": output.Result.Id, "scheduleDate": output.Result.ScheduleTime.String()})
}

func (controller *ScheduleController) DeactivateSchedule(ctx *gin.Context) {
	scheduleId := ctx.Param("scheduleId")
	accountId := ctx.GetString("AccountId")

	output := controller.deactivateWeatherNotification.Execute(notification_schedule_services.DeactivateWeatherNotificationScheduleInputDTO{
		AccountId:  accountId,
		ScheduleId: scheduleId,
	})

	if output.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": output.Error.Message, "code": output.Error.Name})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": output.Result.ScheduleId})
}

func (controller *ScheduleController) ActivateSchedule(ctx *gin.Context) {
	scheduleId := ctx.Param("scheduleId")
	accountId := ctx.GetString("AccountId")

	output := controller.activateWeatherNotification.Execute(notification_schedule_services.ActivateWeatherNotificationScheduleInputDTO{
		AccountId:  accountId,
		ScheduleId: scheduleId,
	})

	if output.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": output.Error.Message, "code": output.Error.Name})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": output.Result.ScheduleId})
}

func (controller *ScheduleController) ListAccountSchedules(ctx *gin.Context) {
	accountId := ctx.GetString("AccountId")

	output := controller.listAccountWeatherNotificationSchedules.Execute(notification_schedule_query.ListAccountWeatherNotificationsServiceInputDTO{
		AccountId: accountId,
	})

	if output.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": output.Error.Message, "code": output.Error.Name})
		return
	}

	response := make([]ListAccountSchedulesResponse, 0, len(*output.Result))

	for _, scheduledNotification := range *output.Result {
		response = append(response, ListAccountSchedulesResponse{
			Id:             scheduledNotification.Id,
			ScheduledDate:  scheduledNotification.ScheduledDate,
			IntervalInDays: scheduledNotification.IntervalInDays,
			Hour:           scheduledNotification.Hour,
			Active:         scheduledNotification.Active,
			City: ListAccountSchedulesCityResponse{
				Id:        scheduledNotification.City.Id,
				Name:      scheduledNotification.City.Name,
				StateCode: scheduledNotification.City.StateCode,
				IsCoastal: scheduledNotification.City.IsCoastal,
			},
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"scheduledNotifications": response})
}

type ScheduleControllerRequest struct {
	Hour           uint8  `json:"hour"`
	IntervalInDays uint8  `json:"intervalInDays"`
	CityId         string `json:"cityId"`
	Method         string `json:"method"`
	IsCoastalCity  bool   `json:"isCoastalCity"`
}

type ListAccountSchedulesResponse struct {
	Id             string                           `json:"id"`
	ScheduledDate  time.Time                        `json:"scheduledDate"`
	IntervalInDays uint8                            `json:"intervalInDays"`
	Hour           uint8                            `json:"hour"`
	Active         bool                             `json:"active"`
	City           ListAccountSchedulesCityResponse `json:"city"`
}

type ListAccountSchedulesCityResponse struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	StateCode string `json:"stateCode"`
	IsCoastal bool   `json:"isCoastal"`
}

func NewScheduleController(
	scheduleWeatherNotification notification_schedule_services.ScheduleNotificationApplicationService,
	deactivateWeatherNotification notification_schedule_services.DeactivateWeatherNotificationScheduleService,
	activateWeatherNotification notification_schedule_services.ActivateWeatherNotificationScheduleService,
	listAccountWeatherNotificationSchedules notification_schedule_query.ListAccountWeatherNotificationsService,
) *ScheduleController {
	return &ScheduleController{
		scheduleWeatherNotification:             scheduleWeatherNotification,
		deactivateWeatherNotification:           deactivateWeatherNotification,
		activateWeatherNotification:             activateWeatherNotification,
		listAccountWeatherNotificationSchedules: listAccountWeatherNotificationSchedules,
	}
}
