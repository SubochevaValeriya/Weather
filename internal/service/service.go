package service

import (
	"github.com/SubochevaValeriya/microservice-balance"
	"github.com/SubochevaValeriya/microservice-balance/internal/repository"
)

type Weather interface {
	AddCity(city string) error
	GetSubscriptionList() ([]weather.Subscription, error)
	GetAvgTempByCity(city string) (int, error)
	DeleteCity(city string) error
	AddWeather(city string) error
}

type Service struct {
	Weather
}

func NewService(repos *repository.Repository) *Service {
	return &Service{newApiService(repos.Balance)}
}
