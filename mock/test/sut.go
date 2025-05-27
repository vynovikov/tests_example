package test

import (
	"context"
	"log/slog"
	"testing"

	apiprivate "somerepo/api/api_private"
	apipublic "somerepo/api/api_public"
	"somerepo/config"
	usecasesUsers "somerepo/usecases/users"
	usecasesUsersMocks "somerepo/usecases/users/mocks"

	"github.com/golang/mock/gomock"
)

type Sut struct {
	MockRepository     *usecasesUsersMocks.MockStorage
	MockUserRepository *usecasesUsersMocks.MockUserRepository
	Ctrl               *gomock.Controller
}

func NewPrivateSut(t *testing.T) (apiprivate.API, Sut) {
	t.Helper()
	ctrl := gomock.NewController(t)

	cfg, err := config.Parse()
	if err != nil {
		t.Fatal(err)
	}

	mockStorageRepository := usecasesUsersMocks.NewMockStorage(ctrl)
	mockUserRepository := usecasesUsersMocks.NewMockUserRepository(ctrl)
	usecases := usecasesUsers.New(mockStorageRepository, mockUserRepository, slog.Default())
	api := apiprivate.New(cfg.ServiceID, usecases, slog.Default())

	return api, Sut{
		MockRepository:     mockStorageRepository,
		MockUserRepository: mockUserRepository,
		Ctrl:               ctrl,
	}
}

func NewPublicSut(t *testing.T) (apipublic.API, Sut, context.Context) {
	t.Helper()
	ctrl := gomock.NewController(t)

	cfg, err := config.Parse()
	if err != nil {
		t.Fatal(err)
	}

	mockStorageRepository := usecasesUsersMocks.NewMockStorage(ctrl)
	mockUserRepository := usecasesUsersMocks.NewMockUserRepository(ctrl)
	usecases := usecasesUsers.New(mockStorageRepository, mockUserRepository, slog.Default())
	api := apipublic.New(cfg.ServiceID, usecases, slog.Default())

	return api, Sut{
		MockRepository:     mockStorageRepository,
		MockUserRepository: mockUserRepository,
		Ctrl:               ctrl,
	}, t.Context()
}
