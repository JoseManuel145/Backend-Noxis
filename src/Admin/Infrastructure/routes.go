package Infrastructure

import (
	"Backend/src/Admin/Infrastructure/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	loginController *handlers.LoginController,
	logoutController *handlers.LogOutController,
	registerController *handlers.RegisterController,
) {
	user := router.Group("/user")
	{
		user.POST("/logIn", loginController.Run)
		user.POST("/logout", logoutController.Run)
		user.POST("/register", registerController.Run)
	}
}
