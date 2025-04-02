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
	clave := ctx.Param("clave")

	id, err := strconv.Atoi(userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "El userId debe ser un número entero"})
		return
	}

	err = uc.usecase.Execute(clave, true, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Kit canjeado correctamente"})
}
