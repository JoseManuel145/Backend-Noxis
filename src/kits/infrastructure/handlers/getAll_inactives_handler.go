package handlers

import (
	"Backend/src/Kits/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetInactivesKits struct {
	usecase *application.GetAllInactives
}

func NewGetInactivesKits(uc *application.GetAllInactives) *GetInactivesKits {
	return &GetInactivesKits{usecase: uc}
}
func (uc *GetInactivesKits) Run(ctx *gin.Context) {
	kits, err := uc.usecase.Execute()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, kits)
}
