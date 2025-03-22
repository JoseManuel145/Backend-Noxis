package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type LogOutController struct{}

func NewLogOutController() *LogOutController {
	return &LogOutController{}
}

func (logout *LogOutController) Run(c *gin.Context) {
	// Eliminar la cookie
	c.SetCookie(
		"jwt",       // name
		"",          // value
		0,           // expires (0 para que se elimine)
		"/",         // path
		"localhost", // domain
		false,       // secure (false para desarrollo)
		true,        // httpOnly (para prevenir XSS)
	)

	// Respuesta
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
