package notification_domain_entities

import (
	"fmt"
	"time"
	"weather_notification/src/modules/shared/value_objects"
)

type WeatherNotificationSchedule struct {
	Id             value_objects.ID
	City           City
	IntervalInDays int8
	Active         bool
	Method         Method
	Hour           Hour
	ScheduledDate  time.Time
}

type City struct {
	Id        string
	IsCoastal bool
}

func (entity WeatherNotificationSchedule) GetId() string {
	return entity.Id.Value
}

func (entity *WeatherNotificationSchedule) Deactivate() {
	entity.Active = false
}

func (entity *WeatherNotificationSchedule) Activate() {
	entity.Active = true
}

func (entity *WeatherNotificationSchedule) calculateScheduleDate() time.Time {
	scheduleDate := entity.ScheduledDate
	scheduleDate = time.Date(scheduleDate.Year(), scheduleDate.Month(), scheduleDate.Day(), 0, 0, 0, 0, scheduleDate.Location())
	scheduleDate = scheduleDate.Add(time.Duration(entity.Hour.Value) * time.Hour)
	scheduleDate = scheduleDate.AddDate(0, 0, int(entity.IntervalInDays))

	return scheduleDate
}

func NewWeatherNotificationSchedule(command NewWeatherNotificationScheduleCommand) (*WeatherNotificationSchedule, error) {
	hour, invalidHour := NewHour(command.Hour)

	if invalidHour != nil {
		return nil, invalidHour
	}

	method, invalidMethod := NewMethod(command.Method)

	if invalidMethod != nil {
		return nil, invalidMethod
	}

	result := &WeatherNotificationSchedule{
		Id:             value_objects.NewID(),
		City:           City{Id: command.CityId, IsCoastal: command.IsCityCoastal},
		IntervalInDays: command.IntervalInDays,
		Active:         true,
		ScheduledDate:  time.Now(),
		Method:         *method,
		Hour:           *hour,
	}

	result.ScheduledDate = result.calculateScheduleDate()

	return result, nil
}

func RecoverWeatherNotificationSchedule(command RecoverWeatherNotificationScheduleCommand) WeatherNotificationSchedule {
	return WeatherNotificationSchedule{
		Id:             value_objects.RecoverID(command.Id),
		City:           City{Id: command.CityId, IsCoastal: command.IsCityCoastal},
		IntervalInDays: command.IntervalInDays,
		Active:         command.Active,
		Method:         Method(command.Method),
		ScheduledDate:  command.ScheduledDate,
		Hour:           Hour{Value: command.Hour},
	}
}

type NewWeatherNotificationScheduleCommand struct {
	CityId         string
	IsCityCoastal  bool
	IntervalInDays int8
	Hour           int8
	Method         string
}

type RecoverWeatherNotificationScheduleCommand struct {
	Id             string
	CityId         string
	IsCityCoastal  bool
	IntervalInDays int8
	Active         bool
	Hour           int8
	Method         string
	ScheduledDate  time.Time
}

type Hour struct {
	Value int8
}

func NewHour(value int8) (*Hour, error) {
	if value < 0 || value > 23 {
		return nil, fmt.Errorf("provide a valid hour between 0 and 23")
	}

	return &Hour{Value: value}, nil
}

type Method string

const (
	Email Method = "EMAIL"
	Sms   Method = "SMS"
	Web   Method = "WEB"
)

func NewMethod(method string) (*Method, error) {
	methodOutput := Method(method)
	switch methodOutput {
	case Email, Sms, Web:
		return &methodOutput, nil
	default:
		return nil, fmt.Errorf("invalid method")
	}
}
