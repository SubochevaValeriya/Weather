package service

import (
	"github.com/SubochevaValeriya/microservice-weather"
	"github.com/SubochevaValeriya/microservice-weather/internal/repository"
	"github.com/SubochevaValeriya/microservice-weather/internal/service/openWeatherApi"
)

type ApiService struct {
	repo        repository.Weather
	openWeather openWeatherApi.Weather
}

func newApiService(repo repository.Weather, openWeather openWeatherApi.Weather) *ApiService {
	return &ApiService{repo: repo, openWeather: openWeather}
}

func (s *ApiService) AddCity(city string) error {
	if _, err := s.openWeather.CurrentTemperature(city); err != nil {
		return err
	}

	if err := s.repo.AddCity(city); err != nil {
		return err
	}

	// add current temperature when adding city
	return s.AddWeather(city)
}

func (s *ApiService) GetSubscriptionList() ([]weather.Subscription, error) {

	return s.repo.GetSubscriptionList()
}

func (s *ApiService) GetAvgTempByCity(city string) (float64, error) {
	if _, err := s.openWeather.CurrentTemperature(city); err != nil {
		return 0, err
	}

	return s.repo.GetAvgTempByCity(city)
}

func (s *ApiService) DeleteCity(city string) error {
	if _, err := s.openWeather.CurrentTemperature(city); err != nil {
		return err
	}

	return s.repo.DeleteCity(city)
}

func (s *ApiService) AddWeather(city string) error {

	temperature, err := s.openWeather.CurrentTemperature(city)
	if err != nil {
		return err
	}

	return s.repo.AddWeather(city, temperature)
}
