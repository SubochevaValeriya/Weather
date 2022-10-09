package service

import (
	"github.com/SubochevaValeriya/microservice-balance"
	"github.com/SubochevaValeriya/microservice-balance/internal/repository"
)

type ApiService struct {
	repo repository.Balance
}

func newApiService(repo repository.Balance) *ApiService {
	return &ApiService{repo: repo}
}

func (s *ApiService) AddCity(city string) error {

	if _, err := CurrentTemperature(city); err != nil {
		return err
	}

	return s.repo.AddCity(city)
}

func (s *ApiService) GetSubscriptionList() ([]weather.Subscription, error) {

	return s.repo.GetSubscriptionList()
}

func (s *ApiService) GetAvgTempByCity(city string) (int, error) {
	if _, err := CurrentTemperature(city); err != nil {
		return 0, err
	}

	return s.repo.GetAvgTempByCity(city)
}

func (s *ApiService) DeleteCity(city string) error {
	if _, err := CurrentTemperature(city); err != nil {
		return err
	}

	return s.repo.DeleteCity(city)
}

func (s *ApiService) AddWeather(city string) error {

	temperature, err := CurrentTemperature(city)
	if err != nil {
		return err
	}

	return s.repo.AddWeather(city, temperature)
}
