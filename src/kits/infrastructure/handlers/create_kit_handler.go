package handlers

import (
	"Backend/src/Kits/application"
	"Backend/src/Kits/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateKit struct {
	usecase *application.CreateKit
}

func NewCreateKit(uc *application.CreateKit) *CreateKit {
	return &CreateKit{usecase: uc}
}
func (c *CreateKit) Run(ctx *gin.Context) {
	var kit domain.Kit
	if err := ctx.ShouldBindJSON(&kit); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}
	err := c.usecase.Execute(&kit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create report: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Report created successfully",
		"report":  kit,
	})
}
