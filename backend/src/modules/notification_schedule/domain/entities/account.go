package notification_domain_entities

import (
	"fmt"
	shared_utils "weather_notification/src/modules/shared/utils"
	"weather_notification/src/modules/shared/value_objects"
)

type AccountAggregateRoot struct {
	Id                     value_objects.ID
	Email                  string
	Phone                  string
	ScheduledNotifications shared_utils.EntityCollection[WeatherNotificationSchedule]
}

func RecoverAccountAggregate(command RecoverAccountAggregateCommand) AccountAggregateRoot {
	scheduledNotifications := make([]WeatherNotificationSchedule, 0, len(command.ScheduledNotifications))

	for _, scheduledNotification := range command.ScheduledNotifications {
		scheduledNotifications = append(scheduledNotifications, RecoverWeatherNotificationSchedule(scheduledNotification))
	}

	return AccountAggregateRoot{
		Id:                     value_objects.RecoverID(command.Id),
		Email:                  command.Email,
		Phone:                  command.Phone,
		ScheduledNotifications: shared_utils.NewCollection(scheduledNotifications),
	}
}

func (agg *AccountAggregateRoot) ScheduleWeatherNotification(input ScheduleWeatherNotificationInput) error {
	if agg.ScheduledNotifications.Length() == 10 {
		return fmt.Errorf("your account has reached the limit for notifications")
	}

	for _, notification := range agg.ScheduledNotifications.GetItems() {
		if notification.City.Id == input.CityId && notification.City.IsCoastal == input.IsCityCoastal {
			return fmt.Errorf("there is a notification already scheduled for this city")
		}
	}

	notification, validationError := NewWeatherNotificationSchedule(NewWeatherNotificationScheduleCommand{
		CityId:         input.CityId,
		IsCityCoastal:  input.IsCityCoastal,
		Hour:           int8(input.Hour),
		Method:         input.Method,
		IntervalInDays: int8(input.IntervalInDays),
	})

	if validationError != nil {
		return validationError
	}

	agg.ScheduledNotifications.Add(*notification)

	return nil
}

func (agg *AccountAggregateRoot) DeactivateSchedule(scheduleId string) error {
	schedule := agg.ScheduledNotifications.Get(scheduleId)

	if schedule == nil {
		return fmt.Errorf("schedule `%s` not found", scheduleId)
	}

	if !schedule.Active {
		return fmt.Errorf("schedule `%s` is already deactivated", scheduleId)
	}

	schedule.Deactivate()

	agg.ScheduledNotifications.Add(*schedule)

	return nil
}

func (agg *AccountAggregateRoot) ActivateSchedule(scheduleId string) error {
	schedule := agg.ScheduledNotifications.Get(scheduleId)

	if schedule == nil {
		return fmt.Errorf("schedule `%s` not found", scheduleId)
	}

	if schedule.Active {
		return fmt.Errorf("schedule `%s` is already activated", scheduleId)
	}

	schedule.Activate()

	agg.ScheduledNotifications.Add(*schedule)

	return nil
}

func (agg AccountAggregateRoot) GetId() string {
	return agg.Id.Value
}

type ScheduleWeatherNotificationInput struct {
	CityId         string
	IsCityCoastal  bool
	IntervalInDays uint8
	Method         string
	Hour           uint8
}

type RecoverAccountAggregateCommand struct {
	Id                     string
	Email                  string
	Phone                  string
	ScheduledNotifications []RecoverWeatherNotificationScheduleCommand
}
