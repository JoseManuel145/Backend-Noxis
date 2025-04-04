package handlers

import (
	"Backend/src/kits/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetActivesKits struct {
	usecase *application.GetAllActives
}

func NewGetActivesKits(uc *application.GetAllActives) *GetActivesKits {
	return &GetActivesKits{usecase: uc}
}
func (uc *GetActivesKits) Run(ctx *gin.Context) {
	kits, err := uc.usecase.Execute()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, kits)
}
