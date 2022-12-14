package handler

import (
	_ "github.com/SubochevaValeriya/microservice-weather/docs"
	"github.com/SubochevaValeriya/microservice-weather/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/weather")
	api.POST("/", h.addCity)              // добавить город в сборщик информации о погоде
	api.GET("/", h.getSubscriptionList)   // получить статистику по списку подписки (сколько и какие города мониторятся)
	api.GET("/:city", h.getAvgTempByCity) // получить среднюю температуру за последние дни (глубина фактически накопленных данных в брокере).
	api.DELETE("/:city", h.deleteCity)    // удалить город из сборщика информации о погоде

	return router
}
