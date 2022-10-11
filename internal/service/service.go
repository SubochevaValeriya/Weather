package service

import (
	"github.com/SubochevaValeriya/microservice-weather"
	"github.com/SubochevaValeriya/microservice-weather/internal/repository"
	"github.com/SubochevaValeriya/microservice-weather/internal/service/openWeatherApi"
)

type Weather interface {
	AddCity(city string) error
	GetSubscriptionList() ([]weather.Subscription, error)
	GetAvgTempByCity(city string) (float64, error)
	DeleteCity(city string) error
	AddWeather(city string) error
}

type Service struct {
	Weather
}

func NewService(repos *repository.Repository, openWeather *openWeatherApi.OpenWeather) *Service {
	return &Service{newApiService(repos.Weather, openWeather.Weather)}
}
