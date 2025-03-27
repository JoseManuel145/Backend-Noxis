package Infrastructure

import (
	"Backend/src/Admin/Infrastructure/adapters"
	"Backend/src/Admin/Infrastructure/handlers"
	"Backend/src/Admin/application"

	"github.com/gin-gonic/gin"
)

func InitUser(router *gin.Engine) {
	println("Usuario Master")
	db := adapters.NewMySQL()
	// Instanciar casos de uso (Use Cases)
	logIn := application.NewLogInUseCase(db)
	register := application.NewRegisterUseCase(db)

	// Instanciar controladores (Handlers)
	logInController := handlers.NewLoginController(logIn)
	registerController := handlers.NewRegisterController(register)
	logOutController := handlers.NewLogOutController()

	// Configurar rutas
	SetupRoutes(router, logInController, logOutController, registerController)

}
