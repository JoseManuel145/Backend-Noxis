package handlers

import (
	"Backend/src/Admin/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	loginUseCase *application.LogInUseCase
}

func NewLoginController(useCase *application.LogInUseCase) *LoginController {
	return &LoginController{
		loginUseCase: useCase,
	}
}

func (login *LoginController) Run(c *gin.Context) {
	var data map[string]string
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}
	token, err := login.loginUseCase.Execute(data["email"], data["password"])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Configurar la cookie con el token
	c.SetCookie(
		"jwt",       //name
		token,       //value
		3600,        //Expires in second 3600s = 1h
		"/",         //path
		"localhost", //domain
		false,       //secure
		true,        //httpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   token,
	})
}
