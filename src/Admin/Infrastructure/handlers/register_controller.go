package handlers

import (
	"Backend/src/Admin/app"
	"Backend/src/Admin/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterController struct {
	userRegister *app.RegisterUseCase
}

func NewRegisterController(useCase *app.RegisterUseCase) *RegisterController {
	return &RegisterController{
		userRegister: useCase,
	}
}
func (register *RegisterController) Run(c *gin.Context) {
	var user domain.Admin
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Llamar al caso de uso de registro
	err := register.userRegister.Execute(user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register admin: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Admin registered successfully",
		"admin":   user,
	})
}
