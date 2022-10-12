package test

import (
	"fmt"
	"github.com/SubochevaValeriya/microservice-weather/internal/handler"
	"github.com/SubochevaValeriya/microservice-weather/internal/repository"
	"github.com/SubochevaValeriya/microservice-weather/internal/service"
	"github.com/SubochevaValeriya/microservice-weather/internal/service/openWeatherApi"
	unitTest "github.com/Valiben/gin_unit_test"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// user and password will need to match running postgres instance
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing congigs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     os.Getenv("host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("failed to inititalize db: %s", err.Error())
	}

	// Check if DB is connected
	if err := db.Ping(); err != nil {
		fmt.Print(err.Error())
		panic(err)
	}

	dbTables := repository.DbTables{SubscriptionTable: viper.GetString("dbTables.subscription"),
		WeatherTable:        viper.GetString("dbTables.weather"),
		WeatherArchiveTable: viper.GetString("dbTables.weather_archive")}

	openWeather, err := openWeatherApi.NewOpenWeatherApiConnection(openWeatherApi.Config{
		Unit:  viper.GetString("openWeatherMap.unit"),
		Lang:  viper.GetString("openWeatherMap.lang"),
		Token: os.Getenv("OWM_API_KEY"),
	})

	// dependency injection
	repos := repository.NewRepository(db, dbTables)
	openWeatherAPI := openWeatherApi.NewOpenWeather(openWeather)
	services := service.NewService(repos, openWeatherAPI)
	handlers := handler.NewHandler(services)
	//srv := new(weather.Server)

	router := handlers.InitRoutes()
	unitTest.SetRouter(router)
	newLog := log.New(os.Stdout, "", log.Llongfile|log.Ldate|log.Ltime)
	unitTest.SetLog(newLog)

	createSchema(db)
	addDataToDatabase(db)
	// replaced package DB to our mock DB
	//postgres.DB = db

	log.Println("Database setup for test")
	exitVal := m.Run()
	log.Println("Database dropped after test")

	deleteSchema(db)

	os.Exit(exitVal)
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func createSchema(db *sqlx.DB) error {
	var schema = `
CREATE TABLE subscription_test
(
    id serial not null primary key unique,
    city varchar(255) not null unique,
    subscription_date date not null
);

CREATE TABLE weather_test

(
    id serial not null unique primary key,
    city_id int not null references subscription_test(id),
    weather_date date not null,
    weather float
);

CREATE TABLE weather_archive_test

(
    id serial not null unique primary key,
    city_id int not null references subscription_test(id),
    weather_date date not null,
    weather float
);`

	db.MustExec(schema)
	return nil
}

func deleteSchema(db *sqlx.DB) error {
	schema := `DROP TABLE weather_test;

	DROP TABLE weather_archive_test;

	DROP TABLE subscription_test;`

	db.MustExec(schema)
	return nil
}

var RandomDate = time.Date(2022, 10, 11, 0, 0, 0, 0, time.UTC)

func addDataToDatabase(db *sqlx.DB) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	addToSubscription := fmt.Sprintf("INSERT INTO subscription_test (city, subscription_date) values ($1, $2)")
	_, err = tx.Exec(addToSubscription, "Moscow", RandomDate)
	if err != nil {
		tx.Rollback()
		return err
	}

	addToWeather := fmt.Sprintf("INSERT INTO weather_test (city_id, weather_date, weather) values ($1, $2, $3), ($4, $5, $6)")
	_, err = tx.Exec(addToWeather, 1, RandomDate, 10, 1, RandomDate, 20)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
