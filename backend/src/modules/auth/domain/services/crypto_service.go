package crypto_service

type CryptoService interface {
	Encrypt(val string) (string, error)
	Verify(password, hash string) bool
}
