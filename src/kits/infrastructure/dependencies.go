package infrastructure

import (
	"Backend/src/Kits/application"
	"Backend/src/kits/infrastructure/adapters"
	"Backend/src/kits/infrastructure/handlers"

	"github.com/gin-gonic/gin"
)

func InitKits(router *gin.Engine) {
	println("KITS")

	db := adapters.NewMySQL()
	create := application.NewCreateKit(db)
	getAll := application.NewGetAllKits(db)
	update := application.NewUpdateKit(db)

	createController := handlers.NewCreateKit(create)
	getAllController := handlers.NewGetAllKits(getAll)
	updateController := handlers.NewUpdateKit(update)

	SetupRoutes(router, createController, getAllController, updateController)
}
