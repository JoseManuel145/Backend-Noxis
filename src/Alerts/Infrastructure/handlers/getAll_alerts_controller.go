package handlers

import (
	"Backend/src/Alerts/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetAllController struct {
	getAllAlertsUseCase *application.GetAllAlerts
}

func NewGetAllAlerts(getAllAlertsUseCase *application.GetAllAlerts) *GetAllController {
	return &GetAllController{
		getAllAlertsUseCase: getAllAlertsUseCase,
	}
}
func (c *GetAllController) Run(ctx *gin.Context) {
	alerts, err := c.getAllAlertsUseCase.Execute()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, alerts)
}
