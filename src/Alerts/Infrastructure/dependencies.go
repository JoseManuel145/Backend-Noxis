package Infrastructure

import (
	"Backend/src/Alerts/Infrastructure/adapters"
	"Backend/src/Alerts/application"

	"github.com/gin-gonic/gin"
)

// Dependencies almacena las instancias de los servicios
type Dependencies struct {
	ProcessReportUseCase *application.ProcessSensor
}

// NewDependencies configura las dependencias del sistema
func NewDependencies(router *gin.Engine) {
	// Inicializar servicio RabbitMQ
	rabbitService := adapters.NewRabbitMQAdapter()

	application.NewProcessSensor(rabbitService)

	// Registrar rutas
	//RegisterRoutes(router, processReportUseCase)

	// Iniciar la escucha de reportes pendientes en un goroutine
	rabbitService.FetchReports()

}
