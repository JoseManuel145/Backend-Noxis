package infrastructure

import (
	application "Backend/src/Alerts/app"
	"Backend/src/Alerts/infrastructure/adapters"
	"Backend/src/core"

	"github.com/gin-gonic/gin"
)

// Dependencies almacena las instancias de los servicios
type Dependencies struct {
	ProcessReportUseCase *application.ProcessSensor
}

// NewDependencies configura las dependencias del sistema
func NewDependencies(router *gin.Engine, rabbitConn *core.ConnRabbitMQ) error {
	// Inicializar servicio RabbitMQ
	rabbitService := adapters.NewRabbitMQAdapter(rabbitConn)

	application.NewProcessSensor(rabbitService)

	// Registrar rutas
	//RegisterRoutes(router, processReportUseCase)

	// Iniciar la escucha de reportes pendientes en un goroutine
	go rabbitService.FetchReports()

	return nil
}
