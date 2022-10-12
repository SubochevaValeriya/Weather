package handler

import (
	weather "github.com/SubochevaValeriya/microservice-weather"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

type CityResponse struct {
	City string `json:"city"`
}

type GetSubscriptionListResponse struct {
	Data []weather.Subscription `json:"data"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, ErrorResponse{message})
}

func newSuccessResponse(method string, city string) {
	if city == "" {
		logrus.Printf("Succesful request for %s", method)
	} else {
		logrus.Printf("Succesful request for %s - %s", method, city)
	}
}
