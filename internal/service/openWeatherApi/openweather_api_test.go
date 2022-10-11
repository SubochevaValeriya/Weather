package openWeatherApi

import (
	"fmt"
	owm "github.com/briandowns/openweathermap"
	"github.com/joho/godotenv"
	"os"
	"testing"
)

type TestCase struct {
	Name          string
	city          string
	ExpectedError error
}

type currentWeather struct {
	query   string
	weather owm.CurrentWeatherData
}

func TestCurrentTemperature(t *testing.T) {
	t.Parallel()

	testCities := []currentWeather{
		{
			query: "Moscow",
			weather: owm.CurrentWeatherData{
				Name: "Moscow",
				Main: owm.Main{
					Temp: 2.91,
				},
			},
		},
		{
			query: "Oludeniz",
			weather: owm.CurrentWeatherData{
				Name: "Oludeniz",
				Main: owm.Main{
					Temp: 16.69,
				},
			},
		},
	}

	testBadCities := []string{"nowhere_", "somewhere_over_the_"}
	err := godotenv.Load()
	c, err := owm.NewCurrent("C", "en", os.Getenv("OWM_API_KEY"))

	if err != nil {
		t.Error(err)
	}

	w := NewApiOpenWeather(c)

	for _, city := range testCities {
		w.CurrentTemperature(city.query)

		if c.Name != city.weather.Name {
			t.Errorf("Expect City %s, got %s", city.weather.Name, c.Name)
		}
		fmt.Println(err)
	}

	for _, badCity := range testBadCities {
		if err := c.CurrentByName(badCity); err != nil {
			t.Log("received expected failure for bad city by name")
		}
	}
}
