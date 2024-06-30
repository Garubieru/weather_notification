package auth_application_service

import (
	"testing"
	"weather_notification/src/tests/mocks"

	"github.com/brianvoe/gofakeit/v6"
)

func TestAccountCreation(t *testing.T) {
	memoryAccountRepository := mocks.NewAccountRepositoryInMemory()
	sut := NewCreateAccountApplication(memoryAccountRepository, mocks.NewCryptoServiceMock())

	res := sut.Execute(MockCreateAccountInput())

	if res != nil {
		t.Error("Error", res)
	}
}

func TestUsernameAlreadyTaken(t *testing.T) {
	memoryAccountRepository := mocks.NewAccountRepositoryInMemory()

	sut := NewCreateAccountApplication(memoryAccountRepository, mocks.NewCryptoServiceMock())

	inputA := MockCreateAccountInput()
	inputA.Username = "any"

	resA := sut.Execute(inputA)

	if resA != nil {
		t.Error("Error", resA)
	}

	inputB := MockCreateAccountInput()
	inputB.Username = "any"

	res := sut.Execute(inputB)

	if res.Name != "UsernameAlreadyTaken" {
		t.Errorf("Expected UsernameAlreadyTaken got %s", res.Name)
	}
}

func TestEmailAlreadyTaken(t *testing.T) {
	memoryAccountRepository := mocks.NewAccountRepositoryInMemory()
	sut := NewCreateAccountApplication(memoryAccountRepository, mocks.NewCryptoServiceMock())

	inputA := MockCreateAccountInput()
	inputA.Email = "any"

	resA := sut.Execute(inputA)

	if resA != nil {
		t.Error("Error", resA)
	}

	inputB := MockCreateAccountInput()
	inputB.Email = "any"

	res := sut.Execute(inputB)

	if res.Name != "EmailAlreadyTaken" {
		t.Errorf("Expected EmailAlreadyTaken got %s", res.Name)
	}
}

func MockCreateAccountInput() CreateAccountInput {
	return CreateAccountInput{Name: gofakeit.Name(), Username: gofakeit.Username(), Email: gofakeit.Email(), Password: gofakeit.Word()}
}
