package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Status string `json:"status"`
}

type cityResponse struct {
	City string `json:"city"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}

func newSuccessResponse(method string, city string) {
	if city == "" {
		logrus.Printf("Succesful request for %s", method)
	} else {
		logrus.Printf("Succesful request for %s - %s", method, city)
	}
}
