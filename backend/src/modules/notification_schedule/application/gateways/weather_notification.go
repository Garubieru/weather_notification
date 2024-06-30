package gateways

import "time"

type WeatherGateway interface {
	GetPrediction(days int, cityId string, isCoastal bool) ([]PredictionDTO, error)
}

type PredictionDTO struct {
	Date      time.Time `json:"date"`
	Max       uint8     `json:"max"`
	Min       uint8     `json:"min"`
	Condition string    `json:"condition"`
	Uvi       uint8    `json:"uvi"`
}
