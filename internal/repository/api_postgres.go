package repository

import (
	"fmt"
	"github.com/SubochevaValeriya/microservice-weather"
	"github.com/jmoiron/sqlx"
	"strconv"
	"time"
)

type ApiPostgres struct {
	db *sqlx.DB
}

func NewApiPostgres(db *sqlx.DB) *ApiPostgres {
	return &ApiPostgres{db: db}
}

func (r *ApiPostgres) AddCity(city string, time time.Time) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	addToSubscription := fmt.Sprintf("INSERT INTO %s (city, subscription_date) values ($1, $2)", subscriptionTable)
	_, err = tx.Exec(addToSubscription, city, time)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *ApiPostgres) GetSubscriptionList() ([]weather.Subscription, error) {
	var subscription []weather.Subscription

	query := fmt.Sprintf("SELECT * FROM %s", subscriptionTable)
	err := r.db.Select(&subscription, query)

	return subscription, err
}

func (r *ApiPostgres) AddWeatherByCityId(id int, temperature int) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	addWeather := fmt.Sprintf("INSERT INTO %s (city_id, weather_date, weather) values ($1, $2, $3)", weathersTable)
	_, err = tx.Exec(addWeather, id, time.Now(), temperature)
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

	deleteTransactionsQuery := fmt.Sprintf("DELETE FROM %s WHERE city_id = $1", weathersTable)
	if _, err := r.db.Exec(deleteTransactionsQuery, id); err != nil {
		tx.Rollback()
		return err
	}

	deleteBalanceQuery := fmt.Sprintf("DELETE FROM %s WHERE id = $1", subscriptionTable)
	_, err = tx.Exec(deleteBalanceQuery, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *ApiPostgres) GetCityId(city string) (int, error) {
	var id int
	getCityIdQuery := fmt.Sprintf("SELECT id from %s WHERE city=$1", subscriptionTable)
	err := r.db.Get(&id, getCityIdQuery, city)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (r *ApiPostgres) GetAvgTempByCityId(id int) (float64, error) {
	var avgTemp float64
	getAvgTemp := fmt.Sprintf("SELECT AVG(weather) FROM %s WHERE city_id=$1", weathersTable)

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

func (r *ApiPostgres) MoveOldDataToArchive() error {
	dateForDelete := time.Now()
	moveOldDataToArchiveQuery := fmt.Sprintf("WITH moved_rows AS (DELETE FROM %s WHERE (weather_date) <= $1 RETURNING *) INSERT INTO %s SELECT * FROM moved_rows", weathersTable, weathersTableArchive)
	if _, err := r.db.Exec(moveOldDataToArchiveQuery, dateForDelete); err != nil {
		return err
	}

	return nil
}
