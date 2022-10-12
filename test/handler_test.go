package test

import (
	weather "github.com/SubochevaValeriya/microservice-weather"
	"github.com/SubochevaValeriya/microservice-weather/internal/handler"
	unitTest "github.com/Valiben/gin_unit_test"
	"github.com/Valiben/gin_unit_test/utils"
	"testing"
)

func TestGetAverageTemperatureDefunctCity(t *testing.T) {
	var avgTemp int

	var response handler.ErrorResponse

	err := unitTest.TestHandlerUnMarshalResp(utils.GET, "/weather/notCity", "FORM", nil, &avgTemp)
	if err == nil {
		t.Errorf("Expected error but got %v", avgTemp)
		return
	}

	unitTest.TestHandlerUnMarshalResp(utils.GET, "/weather/notCity", "FORM", nil, &response)

	if response.Message != "defunct city" {
		t.Errorf("Expected defunct city error but got %v", response.Message)
		return
	}

	t.Log("Success")
}

func TestGetAverageTemperature(t *testing.T) {
	t.Parallel()
	var avgTemp int

	err := unitTest.TestHandlerUnMarshalResp(utils.GET, "/weather/Moscow", "FORM", nil, &avgTemp)
	if err != nil {
		t.Errorf("TestPostMappingClientNotFound: %v\n", err)
		return
	}

	if avgTemp != 15 {
		t.Errorf("Expected %v but got %v", 15, avgTemp)
		return
	}

	t.Log("Success")
}

func TestGetAverageTemperatureNoInfo(t *testing.T) {
	var avgTemp int

	var response handler.ErrorResponse

	err := unitTest.TestHandlerUnMarshalResp(utils.GET, "/weather/Oludeniz", "FORM", nil, &avgTemp)
	if err == nil {
		t.Errorf("Expected error but got %v", avgTemp)
		return
	}

	unitTest.TestHandlerUnMarshalResp(utils.GET, "/weather/Oludeniz", "FORM", nil, &response)

	if response.Message != "sql: no rows in result set" {
		t.Errorf("Expected no info error but got %v", response.Message)
		return
	}

	t.Log("Success")
}

func TestGetSubscriptionListIfCityInTheList(t *testing.T) {
	var response handler.GetSubscriptionListResponse

	err := unitTest.TestHandlerUnMarshalResp(utils.GET, "/weather/", "json", nil, &response)
	if err != nil {
		t.Errorf("TestGetActions: %v\n", err)
		return
	}

	if len(response.Data) != 1 {
		t.Errorf("Expected %v but got %v", 1, len(response.Data))
		return
	}

	t.Log("Success")
}

func TestAddNewCity(t *testing.T) {
	t.Parallel()
	var cityResponse handler.CityResponse
	var input = weather.City{
		City: "Oludeniz",
	}

	err := unitTest.TestHandlerUnMarshalResp(utils.POST, "/weather/", "json", input, &cityResponse)
	if err != nil {
		t.Errorf("TestPostMappingClientNotFound: %v\n", err)
		return
	}

	if cityResponse.City != input.City {
		t.Errorf("Expected %v but got %v", input.City, cityResponse.City)
		return
	}

	t.Log("Success")
}

func TestDeleteCity(t *testing.T) {
	t.Parallel()
	var statusResponse handler.StatusResponse

	err := unitTest.TestHandlerUnMarshalResp(utils.DELETE, "/weather/Moscow", "json", nil, &statusResponse)
	if err != nil {
		t.Errorf("TestPostMappingClientNotFound: %v\n", err)
		return
	}

	if statusResponse.Status != "ok" {
		t.Errorf("Expected %v but got %v", "ok", statusResponse.Status)
		return
	}

	t.Log("Success")
}

func TestAddCityThatAlreadyInSubscription(t *testing.T) {
	var cityResponse handler.CityResponse
	var input = weather.City{
		City: "Moscow",
	}

	err := unitTest.TestHandlerUnMarshalResp(utils.POST, "/weather/", "json", input, &cityResponse)
	if err != nil {
		t.Errorf("TestPostMappingClientNotFound: %v\n", err)
		return
	}

	if cityResponse.City == "" {
		var errorResponse handler.ErrorResponse
		unitTest.TestHandlerUnMarshalResp(utils.POST, "/weather/", "json", input, &errorResponse)
		if errorResponse.Message != `pq: duplicate key value violates unique constraint "subscription_test_city_key"` {
			t.Errorf("Expected duplicate key error but got %v", errorResponse.Message)
			return
		}
	}

	t.Log("Success")
}

func TestAddDefunctCity(t *testing.T) {
	t.Parallel()
	var cityResponse handler.CityResponse
	var input = weather.City{
		City: "IAmNotCity",
	}

	err := unitTest.TestHandlerUnMarshalResp(utils.POST, "/weather/", "json", input, &cityResponse)
	if err != nil {
		t.Errorf("TestPostMappingClientNotFound: %v\n", err)
		return
	}

	if cityResponse.City == "" {
		var errorResponse handler.ErrorResponse
		unitTest.TestHandlerUnMarshalResp(utils.POST, "/weather/", "json", input, &errorResponse)
		if errorResponse.Message != `defunct city` {
			t.Errorf("Expected defunct city error but got %v", errorResponse.Message)
			return
		}
	}

	t.Log("Success")
}

func TestDeleteNotCity(t *testing.T) {
	t.Parallel()
	var statusResponse handler.StatusResponse

	err := unitTest.TestHandlerUnMarshalResp(utils.DELETE, "/weather/notcity", "json", nil, &statusResponse)
	if err != nil {
		t.Errorf("TestPostMappingClientNotFound: %v\n", err)
		return
	}

	if statusResponse.Status == "ok" {
		t.Errorf("Expected error but got %v", statusResponse.Status)
		return
	}

	t.Log("Success")
}
