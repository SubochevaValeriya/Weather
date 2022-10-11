package repository

import (
	weather "github.com/SubochevaValeriya/microservice-weather"
	"github.com/jmoiron/sqlx"
)

type Weather interface {
	AddCity(city string) error
	GetSubscriptionList() ([]weather.Subscription, error)
	GetAvgTempByCity(city string) (float64, error)
	DeleteCity(city string) error
	AddWeather(city string, temperature int) error
	MoveOldDataToArchive() error
}

type Repository struct {
	Weather
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{NewApiPostgres(db)}
}
