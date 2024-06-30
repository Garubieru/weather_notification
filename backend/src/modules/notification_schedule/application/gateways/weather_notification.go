package gateways

type WeatherGateway interface {
	GetPrediction(days int, cityId string, isCoastal bool) (*PredictionDTO, error)
}

type PredictionDTO struct {
	Temperatures   []TemperatureDTO  `json:"temperatures"`
	WaveConditions *WaveConditionDTO `json:"waveCondition"`
}

type TemperatureDTO struct {
	Date      string `json:"date"`
	Max       uint8  `json:"max"`
	Min       uint8  `json:"min"`
	Condition string `json:"condition"`
}

type WaveConditionDTO struct {
	Morning   string `json:"morning"`
	Afternoon string `json:"afternoon"`
	Evening   string `json:"evening"`
}
