package repository

import (
	"fmt"
	"github.com/SubochevaValeriya/microservice-weather"
	"github.com/jmoiron/sqlx"
	"strconv"
	"time"
)

type ApiPostgres struct {
	db       *sqlx.DB
	dbTables DbTables
}

type DbTables struct {
	SubscriptionTable   string
	WeatherTable        string
	WeatherArchiveTable string
}

func NewApiPostgres(db *sqlx.DB, dbTables DbTables) *ApiPostgres {
	return &ApiPostgres{db: db,
		dbTables: dbTables}
}

func (r *ApiPostgres) AddCity(city string, date time.Time) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	addToSubscription := fmt.Sprintf("INSERT INTO %s (city, subscription_date) values ($1, $2)", r.dbTables.SubscriptionTable)
	_, err = tx.Exec(addToSubscription, city, date)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *ApiPostgres) GetSubscriptionList() ([]weather.Subscription, error) {
	var subscription []weather.Subscription

	query := fmt.Sprintf("SELECT * FROM %s", r.dbTables.SubscriptionTable)
	err := r.db.Select(&subscription, query)

	return subscription, err
}

func (r *ApiPostgres) AddWeatherByCityId(id int, date time.Time, temperature int) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	addWeather := fmt.Sprintf("INSERT INTO %s (city_id, weather_date, weather) values ($1, $2, $3)", r.dbTables.WeatherTable)
	_, err = tx.Exec(addWeather, id, date, temperature)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *ApiPostgres) DeleteCityById(id int) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	deleteTransactionsQuery := fmt.Sprintf("DELETE FROM %s WHERE city_id = $1", r.dbTables.WeatherTable)
	if _, err := r.db.Exec(deleteTransactionsQuery, id); err != nil {
		tx.Rollback()
		return err
	}

	deleteSubscriptionQuery := fmt.Sprintf("DELETE FROM %s WHERE id = $1", r.dbTables.SubscriptionTable)
	_, err = tx.Exec(deleteSubscriptionQuery, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *ApiPostgres) GetCityId(city string) (int, error) {
	var id int
	getCityIdQuery := fmt.Sprintf("SELECT id from %s WHERE city=$1", r.dbTables.SubscriptionTable)
	err := r.db.Get(&id, getCityIdQuery, city)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (r *ApiPostgres) GetAvgTempByCityId(id int) (float64, error) {
	var avgTemp float64
	getAvgTemp := fmt.Sprintf("SELECT AVG(weather) FROM %s WHERE city_id=$1", r.dbTables.WeatherTable)

	var avgTempUint8 []uint8
	err := r.db.Get(&avgTempUint8, getAvgTemp, id)
	if err != nil {
		return avgTemp, err
	}

	avgTemp, err = strconv.ParseFloat(string(avgTempUint8), 64)
	if err != nil {
		return avgTemp, err
	}
	return avgTemp, nil
}

func (r *ApiPostgres) MoveOldDataToArchive(dateForDelete time.Time) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	moveOldDataToArchiveQuery := fmt.Sprintf("WITH moved_rows AS (DELETE FROM %s WHERE (weather_date) <= $1 RETURNING *) INSERT INTO %s SELECT * FROM moved_rows", r.dbTables.WeatherTable, r.dbTables.WeatherArchiveTable)
	if _, err := r.db.Exec(moveOldDataToArchiveQuery, dateForDelete); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
