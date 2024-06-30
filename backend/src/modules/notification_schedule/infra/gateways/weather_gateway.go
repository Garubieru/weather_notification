package infra_gateways

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"weather_notification/src/modules/notification_schedule/application/gateways"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/html/charset"
)

type CPTECWeatherGateway struct {
	baseUrl     string
	redisClient redis.Client
	ctx         context.Context
}

func NewCPTECWeatherGateway(baseUrl string,
	redisClient redis.Client,
	ctx context.Context) CPTECWeatherGateway {
	return CPTECWeatherGateway{
		baseUrl:     baseUrl,
		redisClient: redisClient,
		ctx:         ctx,
	}
}

func (gateway CPTECWeatherGateway) GetPrediction(days int, cityId string, isCoastal bool) (*gateways.PredictionDTO, error) {
	cacheKey := gateway.getCacheKey(days, cityId, isCoastal)

	var output gateways.PredictionDTO

	cacheValue, err := gateway.redisClient.Get(gateway.ctx, cacheKey).Bytes()

	if err != nil {
		if err != redis.Nil {
			return nil, err
		}
	}

	if len(cacheValue) > 0 {
		err = json.Unmarshal(cacheValue, &output)

		if err != nil {
			return nil, err
		}

		return &output, nil
	}

	data, err := gateway.getPredictionData(cityId, isCoastal)

	if err != nil {
		return nil, err
	}

	var predictionResponse PredictionResponse

	err = xml.Unmarshal(data.temperatures, &predictionResponse)

	if err != nil {
		return nil, err
	}

	temperatures := make([]gateways.TemperatureDTO, 0, len(predictionResponse.Predictions))

	for _, prediction := range predictionResponse.Predictions {
		temperatures = append(temperatures, gateways.TemperatureDTO{
			Date:      prediction.Date,
			Max:       prediction.Max,
			Min:       prediction.Min,
			Condition: prediction.Condition,
		})
	}

	output = gateways.PredictionDTO{
		Temperatures:   temperatures,
		WaveConditions: nil,
	}

	if isCoastal && len(data.waveConditions) > 0 {
		var waveConditionResponse WaveConditionsResponse

		err = xml.Unmarshal(data.waveConditions, &waveConditionResponse)

		if err != nil {
			return nil, err
		}

		if waveConditionResponse.StateCode != "undefined" {
			output.WaveConditions = &gateways.WaveConditionDTO{
				Morning:   waveConditionResponse.Morning.Condition,
				Afternoon: waveConditionResponse.Afternoon.Condition,
				Evening:   waveConditionResponse.Evening.Condition,
			}
		}
	}

	now := time.Now()
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), now.Location())

	predictionJson, err := json.Marshal(output)

	if err != nil {
		return nil, fmt.Errorf("error while parsing struct to json")
	}

	gateway.redisClient.Set(gateway.ctx, cacheKey, predictionJson, time.Until(endOfDay))

	return &output, nil
}

func (gateway CPTECWeatherGateway) getCacheKey(days int, cityId string, isCoastal bool) string {
	return fmt.Sprintf("predictions:%d:%s:%t", days, cityId, isCoastal)
}

func (gateway CPTECWeatherGateway) getPredictionData(cityId string, isCoastal bool) (*getPredictionDataOutput, error) {
	temperatures, predictionRequestErr := gateway.request(gateway.buildUrl(
		[]string{"cidade", cityId, "previsao.xml"},
	))

	if predictionRequestErr != nil {
		return nil, predictionRequestErr
	}

	var waveConditions []byte

	if isCoastal {
		waveConditionResponse, waveConditionRequestErr := gateway.request(gateway.buildUrl(
			[]string{"cidade", cityId, "dia", "0", "ondas.xml"},
		))

		if waveConditionRequestErr != nil {
			return nil, predictionRequestErr
		}

		waveConditions = waveConditionResponse
	}

	return &getPredictionDataOutput{
		temperatures:   temperatures,
		waveConditions: waveConditions,
	}, nil
}

func (gateway CPTECWeatherGateway) request(url string) ([]byte, error) {
	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, NewHttpError(response.StatusCode, "failed to retrieve prediction information")
	}

	contentType := response.Header.Get("Content-Type")

	charsetReader, err := charset.NewReader(response.Body, contentType)

	if err != nil {
		return nil, err
	}

	body, readErr := io.ReadAll(charsetReader)

	if readErr != nil {
		return nil, readErr
	}

	body = gateway.skipXMLDeclaration(body)

	return body, nil
}

func (gateway CPTECWeatherGateway) buildUrl(path []string) string {
	return fmt.Sprintf("%s/%s", gateway.baseUrl, strings.Join(path, "/"))
}

func (gateway CPTECWeatherGateway) skipXMLDeclaration(body []byte) []byte {
	if bytes.HasPrefix(body, []byte("<?xml")) {
		closeBracketIndex := bytes.Index(body, []byte(">"))
		body = body[closeBracketIndex+1:]
	}
	return body
}

type HttpError struct {
	status  int
	message string
}

func NewHttpError(status int, message string) HttpError {
	return HttpError{status: status, message: message}
}

func (httpError HttpError) Error() string {
	return fmt.Errorf("%d - %s", httpError.status, httpError.message).Error()
}

type PredictionResponse struct {
	XMLName     xml.Name     `xml:"cidade"`
	Predictions []Prediction `xml:"previsao"`
}

type Prediction struct {
	Date      string `xml:"dia"`
	Condition string `xml:"tempo"`
	Max       uint8  `xml:"maxima"`
	Min       uint8  `xml:"minima"`
}

type WaveConditionsResponse struct {
	XMLName   xml.Name      `xml:"cidade"`
	StateCode string        `xml:"uf"`
	Morning   WaveCondition `xml:"manha"`
	Afternoon WaveCondition `xml:"tarde"`
	Evening   WaveCondition `xml:"noite"`
}

type WaveCondition struct {
	Condition string `xml:"agitacao"`
}

type getPredictionDataOutput struct {
	temperatures   []byte
	waveConditions []byte
}
