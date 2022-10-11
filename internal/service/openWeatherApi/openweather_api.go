package openWeatherApi

import (
	"errors"
	owm "github.com/briandowns/openweathermap"
)

var defunctCityError = errors.New("defunct city")

type Weather interface {
	CurrentTemperature(city string) (int, error)
}

type OpenWeather struct {
	Weather
}

type ApiOpenWeather struct {
	openWeather *owm.CurrentWeatherData
}

func NewApiOpenWeather(openWeather *owm.CurrentWeatherData) *ApiOpenWeather {
	return &ApiOpenWeather{openWeather: openWeather}
}

func NewOpenWeather(openWeather *owm.CurrentWeatherData) *OpenWeather {
	return &OpenWeather{NewApiOpenWeather(openWeather)}
}

// CurrentTemperature is used to receive current temperature in chosen city
func (w *ApiOpenWeather) CurrentTemperature(city string) (int, error) {
	w.openWeather.CurrentByName(city)

	if w.openWeather.Name == "" {
		return 0.0, defunctCityError
	}

	return int(w.openWeather.Main.Temp), nil
}
