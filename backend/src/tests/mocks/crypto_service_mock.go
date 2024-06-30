package mocks

import (
	services "weather_notification/src/modules/auth/domain/services"

	"github.com/brianvoe/gofakeit/v6"
)

type CryptoServiceMock struct {
	VerifyOutput bool
}

func (service CryptoServiceMock) Encrypt(val string) (string, error) {
	value := gofakeit.Password(true, true, true, true, true, 10)
	return value, nil
}

func (service CryptoServiceMock) Verify(val string, secret string) bool {
	return service.VerifyOutput
}

func NewCryptoServiceMock() services.CryptoService {
	return &CryptoServiceMock{VerifyOutput: true}
}
