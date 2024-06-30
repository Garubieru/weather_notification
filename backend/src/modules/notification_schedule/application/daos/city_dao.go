package notification_schedule_daos

type CityDAO interface {
	FindById(id string) (*CityDTO, error)
}

type CityDTO struct {
	Id         string
	ExternalId string
	Name       string
	StateCode  string
}
