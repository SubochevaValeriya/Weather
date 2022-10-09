package service

import (
	"errors"
	owm "github.com/briandowns/openweathermap"
	"log"
	"os"
)

// не работает
var apiKey = os.Getenv("OWM_API_KEY")
var defunctCityError = errors.New("defunct city")

const (
	unit = "C" // fahrenheit
	lang = "en"
)

func CurrentTemperature(city string) (int, error) {
	w, err := owm.NewCurrent(unit, lang, "60ccce665c20e05f092484ed03193b52")
	if err != nil {
		log.Fatalln(err)
	}

	w.CurrentByName(city)

	if w.Name == "" {
		return 0.0, defunctCityError
	}

	return int(w.Main.Temp), nil
}
