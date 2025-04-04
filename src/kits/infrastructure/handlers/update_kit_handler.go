package handlers

import (
	"Backend/src/Kits/application"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UpdateKit struct {
	usecase *application.UpdateKit
}

func NewUpdateKit(uc *application.UpdateKit) *UpdateKit {
	return &UpdateKit{usecase: uc}
}
func (uc *UpdateKit) Run(ctx *gin.Context) {
	userId := ctx.Param("userId")

	id, err := strconv.Atoi(userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "El userId debe ser un número entero"})
		return
	}
	var requestBody struct {
		Clave string `json:"clave"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al leer el JSON: " + err.Error()})
		return
	}

	err = uc.usecase.Execute(requestBody.Clave, true, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Kit canjeado correctamente"})
}
