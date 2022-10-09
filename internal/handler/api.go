package handler

import (
	weather "github.com/SubochevaValeriya/microservice-balance"
	"github.com/gin-gonic/gin"
	"net/http"
)

// addCity is made for adding city to the subscription
// It's CREATE in CRUD
func (h *Handler) addCity(c *gin.Context) {
	var input weather.Subscription

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Weather.AddCity(input.City)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"city": input.City,
	})

	newSuccessResponse("adding city", input.City)
}

type getSubscriptionListResponse struct {
	Data []weather.Subscription `json:"data"`
}

// getSubscriptionList is used to get list of cities in subscription
// It's READ from CRUD
func (h *Handler) getSubscriptionList(c *gin.Context) {

	subscriptionList, err := h.services.Weather.GetSubscriptionList()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getSubscriptionListResponse{
		Data: subscriptionList,
	})

	newSuccessResponse("getting subscription list", "")
}

// getAvgTempByCity allows to get average temperature in city
// It's READ from CRUD
func (h *Handler) getAvgTempByCity(c *gin.Context) {
	city := c.Param("city")

	temp, err := h.services.Weather.GetAvgTempByCity(city)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, temp)
	newSuccessResponse("getting average temperature", city)
}

// deleteCity allows to delete city from the subscription
// It's DELETE from CRUD
func (h *Handler) deleteCity(c *gin.Context) {
	city := c.Param("city")
	//if err != nil {
	//	newErrorResponse(c, http.StatusBadRequest, "invalid id param")
	//	return
	//}

	err := h.services.Weather.DeleteCity(city)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})

	newSuccessResponse("delete city", city)
}

// addWeather is made for adding weather's information
// It's CREATE in CRUD
func (h *Handler) addWeather(c *gin.Context) {
	var input weather.Subscription

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Weather.AddCity(input.City)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"city": input.City,
	})

	newSuccessResponse("adding weather", input.City)
}
