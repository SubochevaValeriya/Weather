package main

import (
	balance "github.com/SubochevaValeriya/microservice-balance"
	"github.com/SubochevaValeriya/microservice-balance/internal/handler"
	"github.com/SubochevaValeriya/microservice-balance/internal/repository"
	"github.com/SubochevaValeriya/microservice-balance/internal/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"time"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing congigs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	//sudo docker run --name=balance -e POSTGRES_PASSWORD='qwerty' -p 5432:5432 -d --rm postgres
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

	// dependency injection
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	srv := new(balance.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error ccured while running http server: %s", err.Error())
	}

	// periodically check temperature in chosen cities
	timer := time.AfterFunc(timeDurations(viper.GetInt("periodicity")), func() {})

	for range timer.C {
		subscriptionList, err := repos.GetSubscriptionList()
		if err != nil {
			return
		}

		for i := range subscriptionList {
			go services.AddWeather(subscriptionList[i].City)
		}
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func timeDurations(minutes int) time.Duration {
	return time.Duration(minutes * 1e9 * 60)
}
