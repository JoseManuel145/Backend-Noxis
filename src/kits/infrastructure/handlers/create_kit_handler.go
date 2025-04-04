package handlers

import (
	"Backend/src/kits/application"
	"Backend/src/kits/domain"
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
	var request struct {
		Sensores []string `json:"sensores" binding:"required"`
		Username string   `json:"username" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	kit := domain.Kit{
		Sensores: request.Sensores,
		Username: request.Username,
	}

	err := c.usecase.Execute(&kit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create kit: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Kit created successfully",
		"clave":   kit.Clave,
	})
}
