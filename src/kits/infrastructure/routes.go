package infrastructure

import (
	"Backend/src/Kits/infrastructure/handlers"
	"Backend/src/utils"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	createController *handlers.CreateKit,
	getAllController *handlers.GetAllKits,
	updateController *handlers.UpdateKit,
	getActivesController *handlers.GetActivesKits,
	getInactivesController *handlers.GetInactivesKits,
) {
	kits := router.Group("/kits")
	{
		kits.POST("/create", utils.VerificarToken, createController.Run)
		kits.GET("/", utils.VerificarToken, getAllController.Run)
		kits.POST("/:userId", utils.VerificarToken, updateController.Run)
		kits.GET("/actives", utils.VerificarToken, getActivesController.Run)
		kits.GET("/inactives", utils.VerificarToken, getInactivesController.Run)
	}
}
