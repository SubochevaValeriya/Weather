package repository

import (
	weather "github.com/SubochevaValeriya/microservice-weather"
	"github.com/jmoiron/sqlx"
	"time"
)

type Weather interface {
	AddCity(city string, time time.Time) error
	GetSubscriptionList() ([]weather.Subscription, error)
	AddWeatherByCityId(id int, date time.Time, temperature int) error
	MoveOldDataToArchive(dateForDelete time.Time) error
	GetAvgTempByCityId(id int) (float64, error)
	GetCityId(city string) (int, error)
	DeleteCityById(id int) error
}

type Repository struct {
	Weather
}

func NewRepository(db *sqlx.DB, dbTables DbTables) *Repository {
	return &Repository{NewApiPostgres(db, dbTables)}
}
