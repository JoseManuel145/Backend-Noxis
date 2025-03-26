package Infrastructure

import (
	"Backend/src/Admin/Infrastructure/adapters"
	"Backend/src/Admin/Infrastructure/handlers"
	"Backend/src/Admin/app"

	"github.com/gin-gonic/gin"
)

func InitUser(db *adapters.MySQL, router *gin.Engine) {
	println("Usuario Master")

	// Instanciar casos de uso (Use Cases)
	logIn := app.NewLogInUseCase(db)
	register := app.NewRegisterUseCase(db)

	// Instanciar controladores (Handlers)
	logInController := handlers.NewLoginController(logIn)
	registerController := handlers.NewRegisterController(register)
	logOutController := handlers.NewLogOutController()

	// Configurar rutas
	SetupRoutes(router, logInController, logOutController, registerController)

}
