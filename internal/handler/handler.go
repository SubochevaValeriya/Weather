package handler

import (
	"github.com/SubochevaValeriya/microservice-balance/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	//-	управление списком подписки на погоду по городам
	//-	получить статистику по списку подписки (сколько и какие города мониторятся,
	//	-	добавить город в сборщик информации о погоде
	//-	удалить город из сборщика информации о погоде
	//-	получить среднюю температуру за последние дни (глубина фактически накопленных данных в брокере).

	api := router.Group("/weather")
	api.POST("/", h.addCity)              // добавить город в сборщик информации о погоде
	api.GET("/", h.getSubscriptionList)   // получить статистику по списку подписки (сколько и какие города мониторятся)
	api.GET("/:city", h.getAvgTempByCity) // получить среднюю температуру за последние дни (глубина фактически накопленных данных в брокере).
	api.DELETE("/:city", h.deleteCity)    // удалить город из сборщика информации о погоде
	api.POST("/:city", h.addWeather)      // добавить погоду в городе в сборщик информации о погоде

	return router
}
