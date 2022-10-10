package main

import (
	"context"
	balance "github.com/SubochevaValeriya/microservice-balance"
	"github.com/SubochevaValeriya/microservice-balance/internal/handler"
	"github.com/SubochevaValeriya/microservice-balance/internal/repository"
	"github.com/SubochevaValeriya/microservice-balance/internal/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
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

	go periodicallyCheckTemperature(repos, services, viper.GetInt("periodicity"))
	go moveOldDataToArchive(repos, viper.GetInt("cntDayArchive"))

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

// periodically checking temperature in chosen cities
func periodicallyCheckTemperature(repos *repository.Repository, services *service.Service, periodicity int) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	for range time.Tick(timeDurations(periodicity)) {
		subscriptionList, err := repos.GetSubscriptionList()
		if err != nil {
			return
		}
		wg2 := sync.WaitGroup{}
		wg2.Add(len(subscriptionList))
		for i := range subscriptionList {
			go services.AddWeather(subscriptionList[i].City)
			wg2.Done()
		}
		wg2.Wait()
	}
	wg.Wait()
}

// moves old data to archive table
func moveOldDataToArchive(repos *repository.Repository, cntDayArchive int) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	for range time.Tick(timeDurations(cntDayArchive / 1440)) {
		repos.MoveOldDataToArchive(cntDayArchive)
	}
	wg.Wait()
}

func timeDurations(minutes int) time.Duration {
	return time.Duration(minutes * 1e9 * 60)
}
