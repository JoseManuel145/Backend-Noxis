package infrastructure

import (
	"Backend/src/kits/infrastructure/handlers"
	"Backend/src/utils"

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
		kits.POST("/create", utils.VerificarToken, createController.Run)
		kits.GET("/", utils.VerificarToken, getAllController.Run)
		kits.POST("/:userId/:clave", utils.VerificarToken, updateController.Run)
	}
}
