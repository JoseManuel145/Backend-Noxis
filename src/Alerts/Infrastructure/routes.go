package Infrastructure

import (
	"Backend/src/Alerts/Infrastructure/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	getAll *handlers.GetAllController,
	getBySensor *handlers.GetBySensorController,
) {
	alerts := router.Group("/alerts")
	{
		alerts.GET("/", getAll.Run)
		alerts.GET("/:sensor", getBySensor.Run)
	}
}
