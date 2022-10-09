package repository

import (
	"github.com/SubochevaValeriya/microservice-balance"
	"github.com/jmoiron/sqlx"
)

type Balance interface {
	AddCity(city string) error
	GetSubscriptionList() ([]weather.Subscription, error)
	GetAvgTempByCity(city string) (int, error)
	DeleteCity(city string) error
	AddWeather(city string, temperature int) error
}

type Repository struct {
	Balance
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{NewApiPostgres(db)}
}
