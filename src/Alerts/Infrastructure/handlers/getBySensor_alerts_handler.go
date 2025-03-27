package handlers

import (
	"Backend/src/Alerts/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetBySensorController struct {
	getBySensorUseCase *application.GetBySensorAlerts
}

func NewGetBySensor(getBySensorUseCase *application.GetBySensorAlerts) *GetBySensorController {
	return &GetBySensorController{getBySensorUseCase: getBySensorUseCase}
}

func (c *GetBySensorController) Run(ctx *gin.Context) {
	sensor := ctx.Param("sensor")
	if sensor == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro sensor es obligatorio"})
		return
	}

	alerts, err := c.getBySensorUseCase.Execute(sensor)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, alerts)
}
