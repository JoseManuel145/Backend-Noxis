package infrastructure

import (
	"Backend/src/Kits/application"
	"Backend/src/Kits/infrastructure/adapters"
	"Backend/src/Kits/infrastructure/handlers"

	"github.com/gin-gonic/gin"
)

func InitKits(router *gin.Engine) {
	println("KITS")

	db := adapters.NewMySQL()
	create := application.NewCreateKit(db)
	getAll := application.NewGetAllKits(db)
	update := application.NewUpdateKit(db)
	getInactives := application.NewGetAllInactives(db)
	getActives := application.NewGetAllActives(db)

	createController := handlers.NewCreateKit(create)
	getAllController := handlers.NewGetAllKits(getAll)
	updateController := handlers.NewUpdateKit(update)
	getActivesController := handlers.NewGetActivesKits(getActives)
	getInactivesController := handlers.NewGetInactivesKits(getInactives)

	SetupRoutes(router, createController, getAllController, updateController, getActivesController, getInactivesController)
}
