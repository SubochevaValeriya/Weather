package openWeatherApi

import (
	owm "github.com/briandowns/openweathermap"
)

type Config struct {
	Unit  string
	Lang  string
	Token string
}

// NewOpenWeatherApiConnection is used to connect to NewOpenWeatherMap
func NewOpenWeatherApiConnection(config Config) (*owm.CurrentWeatherData, error) {
	w, err := owm.NewCurrent(config.Unit, config.Lang, config.Token)
	if err != nil {
		return nil, err
	}

	return w, nil
}
