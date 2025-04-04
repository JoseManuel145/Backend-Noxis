package handlers

import (
	"Backend/src/Kits/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetAllKits struct {
	usecase *application.GetAllKits
}

func NewGetAllKits(uc *application.GetAllKits) *GetAllKits {
	return &GetAllKits{usecase: uc}
}
func (uc *GetAllKits) Run(ctx *gin.Context) {
	kits, err := uc.usecase.Execute()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, kits)
}
