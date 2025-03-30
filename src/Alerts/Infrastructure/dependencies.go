package Infrastructure

import (
	"Backend/src/Alerts/Infrastructure/adapters"
	"Backend/src/Alerts/Infrastructure/handlers"
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
	db := adapters.NewMongoAlertRepository()
	websocketService := adapters.NewWebSocketAdapter()
	save := application.NewSaveAlert(db)

	getAll := application.NewGetAllAlerts(db)
	getBySensor := application.NewGetBySensorAlert(db)

	ws := application.NewSendAlertUseCase(websocketService)

	application.NewProcessSensor(rabbitService, save, ws)

	getAllController := handlers.NewGetAllAlerts(getAll)
	getBySensorController := handlers.NewGetBySensor(getBySensor)

	// Registrar rutas
	SetupRoutes(router, getAllController, getBySensorController)

	// Iniciar la escucha de reportes pendientes en un goroutine
	rabbitService.FetchReports()

}
