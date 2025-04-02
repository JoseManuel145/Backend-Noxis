package infrastructure

import (
	"Backend/src/kits/infrastructure/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	createController *handlers.CreateKit,
	getAllController *handlers.GetAllKits,
	updateController *handlers.UpdateKit,
) {
	kits := router.Group("/kits")
	{
		kits.POST("/create", createController.Run)
		kits.GET("/", getAllController.Run)
		kits.POST("/:userId/:clave", updateController.Run)
	}
}
