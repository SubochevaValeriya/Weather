package service

import (
	"github.com/SubochevaValeriya/microservice-weather/internal/repository"
	"github.com/sirupsen/logrus"
	"strconv"
	"sync"
	"time"
)

// PeriodicallyCheckTemperature periodically checking temperature in chosen cities
func PeriodicallyCheckTemperature(repos *repository.Repository, services *Service, periodicity int) {
	period, err := time.ParseDuration(strconv.Itoa(periodicity) + "m")
	if err != nil {
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	for range time.Tick(period) {
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
		logrus.Printf("Added new data for cities")
	}
	wg.Wait()
}

// MoveOldDataToArchive moves old data to archive table
func MoveOldDataToArchive(repos *repository.Repository, cntDayArchive int) {
	period, err := time.ParseDuration(strconv.Itoa(cntDayArchive*24) + "h")
	if err != nil {
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	for range time.Tick(period) {
		repos.MoveOldDataToArchive()
		logrus.Printf("Old data is archived")
	}
	wg.Wait()
}
