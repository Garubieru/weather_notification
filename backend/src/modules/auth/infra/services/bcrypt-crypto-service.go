package infra_services

import (
	crypto_service "weather_notification/src/modules/auth/domain/services"

	"golang.org/x/crypto/bcrypt"
)

type BcryptCryptoService struct{}

func (service BcryptCryptoService) Encrypt(value string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (service BcryptCryptoService) Verify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func NewBcryptCryptoService() crypto_service.CryptoService {
	return BcryptCryptoService{}
}
