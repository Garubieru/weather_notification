package infra_gateways

import (
	"time"
	"weather_notification/src/modules/notification_schedule/application/gateways"
)

type CPTECWeatherGateway struct{}

func (gateway CPTECWeatherGateway) GetPrediction(days int, cityId string, isCoastal bool) ([]gateways.PredictionDTO, error) {
	mock := gateways.PredictionDTO{
		Date:      time.Now(),
		Max:       1,
		Min:       20,
		Condition: "pc",
		Uvi:       0,
	}

	return []gateways.PredictionDTO{mock}, nil
}
