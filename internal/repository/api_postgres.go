package repository

import (
	"fmt"
	"github.com/SubochevaValeriya/microservice-balance"
	"github.com/jmoiron/sqlx"
	"time"
)

type ApiPostgres struct {
	db *sqlx.DB
}

func NewApiPostgres(db *sqlx.DB) *ApiPostgres {
	return &ApiPostgres{db: db}
}

func (r *ApiPostgres) AddCity(city string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	addToSubscription := fmt.Sprintf("INSERT INTO %s (city, subscription_date) values ($1, $2)", subscriptionTable)
	_, err = tx.Exec(addToSubscription, city, time.Now())
	if err != nil {
		tx.Rollback()
		return err
	}

	//addTransactionQuery := fmt.Sprintf("INSERT INTO %s (user_id, amount, reason, transaction_date, transfer_id) values ($1, $2, $3, $4, $5)", weathersTable)
	//_, err = tx.Exec(addTransactionQuery, id, user.Balance, ReasonOpening, time.Now(), id)
	//if err != nil {
	//	tx.Rollback()
	//	return 0, err
	//}

	return tx.Commit()
}

func (r *ApiPostgres) GetSubscriptionList() ([]weather.Subscription, error) {
	var subscription []weather.Subscription

	query := fmt.Sprintf("SELECT * FROM %s", subscriptionTable)
	err := r.db.Select(&subscription, query)

	return subscription, err
}

func (r *ApiPostgres) GetAvgTempByCity(city string) (int, error) {
	var avgTemp int
	getAvgTemp := fmt.Sprintf("SELECT AVG(weather) FROM %s WHERE city=$1", weathersTable)
	err := r.db.Get(&avgTemp, getAvgTemp, city)
	if err != nil {
		return avgTemp, err
	}

	return avgTemp, nil
}

func (r *ApiPostgres) DeleteCity(city string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	deleteTransactionsQuery := fmt.Sprintf("DELETE FROM %s WHERE city = $1", weathersTable)
	if _, err := r.db.Exec(deleteTransactionsQuery, city); err != nil {
		tx.Rollback()
		return err
	}

	deleteBalanceQuery := fmt.Sprintf("DELETE FROM %s WHERE city = $1", subscriptionTable)
	_, err = tx.Exec(deleteBalanceQuery, city)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *ApiPostgres) AddWeather(city string, temperature int) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	addWeather := fmt.Sprintf("INSERT INTO %s (city, weather_date, weather) values ($1, $2, $3) RETURNING id", weathersTable)
	_, err = tx.Exec(addWeather, city, time.Now(), temperature)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
