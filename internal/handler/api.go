package handler

import (
	weather "github.com/SubochevaValeriya/microservice-weather"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// @Summary Add City
// @Tags weather
// @Description add city to the subscription
// @ID add-city
// @Accept  json
// @Produce  json
// @Param input body weather.City true "city name"
// @Success 200 {object} cityResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /weather [post]
// addCity is made for adding city to the subscription
// It's CREATE in CRUD
func (h *Handler) addCity(c *gin.Context) {
	var input weather.City

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Weather.AddCity(uppercaseFirstLetter(input.City))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, cityResponse{
		City: uppercaseFirstLetter(input.City),
	})

	newSuccessResponse("adding city", input.City)
}

type getSubscriptionListResponse struct {
	Data []weather.Subscription `json:"data"`
}

// @Summary Get Subscription List
// @Tags weather
// @Description get list of cities in subscription
// @ID get-subscription-list
// @Accept  json
// @Produce  json
// @Success 200 {object} getSubscriptionListResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /weather [get]
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

// @Summary Get Avg Temperature By City
// @Tags weather
// @Description get average temperature in city
// @ID get-avg-temp-by-city
// @Accept  json
// @Produce  json
// @Param city path string true "city name"
// @Success 200 {integer} integer
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /weather/{city} [get]
// getAvgTempByCity allows to get average temperature in city
// It's READ from CRUD
func (h *Handler) getAvgTempByCity(c *gin.Context) {
	city := uppercaseFirstLetter(c.Param("city"))

	temp, err := h.services.Weather.GetAvgTempByCity(city)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, temp)
	newSuccessResponse("getting average temperature", city)
}

// @Summary Delete City
// @Tags weather
// @Description delete city from the subscription
// @ID delete-city
// @Accept  json
// @Produce  json
// @Param city path string true "city name"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /weather/{city} [delete]
// deleteCity allows to delete city from the subscription
// It's DELETE from CRUD
func (h *Handler) deleteCity(c *gin.Context) {
	city := uppercaseFirstLetter(c.Param("city"))

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

func uppercaseFirstLetter(word string) string {
	return strings.TrimSpace(strings.Title(strings.ToLower(word)))
}
