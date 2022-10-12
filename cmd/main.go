package main

import (
	"context"
	weather "github.com/SubochevaValeriya/microservice-weather"
	"github.com/SubochevaValeriya/microservice-weather/internal/handler"
	"github.com/SubochevaValeriya/microservice-weather/internal/repository"
	"github.com/SubochevaValeriya/microservice-weather/internal/service"
	"github.com/SubochevaValeriya/microservice-weather/internal/service/openWeatherApi"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

// @title Weather App API
// @version 1.0
// @description API Server for Weather Application

// @host localhost:8000
// @BasePath /

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing congigs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	//connect w/t docker-compose:
	//sudo docker run --name=weather -e POSTGRES_PASSWORD='qwerty' -p 5432:5432 -d --rm postgres
	// migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable' up

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
	srv := new(weather.Server)

	go service.PeriodicallyCheckTemperature(repos, services, viper.GetInt("periodicity"))
	go service.MoveOldDataToArchive(repos, viper.GetInt("cntDayArchive"))

	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Print("Weather App Started")

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("Weather App Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
