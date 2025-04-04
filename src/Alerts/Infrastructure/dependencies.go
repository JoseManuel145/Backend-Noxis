package Infrastructure

import (
	"Backend/src/Alerts/Infrastructure/adapters"
	"Backend/src/Alerts/Infrastructure/handlers"
	"Backend/src/Alerts/application"
	"log"

	"github.com/gin-gonic/gin"
)

// InitAlerts configura las dependencias del sistema
func InitAlerts(router *gin.Engine) {
	println("SENSORES")

	// Inicializar servicio RabbitMQ
	rabbitService := adapters.NewRabbitMQAdapter()

	// Inicializar repositorios y servicios
	db := adapters.NewMongoAlertRepository()
	websocketService := adapters.NewWebSocketAdapter()
	save := application.NewSaveAlert(db)

	// Inicializar casos de uso
	getAll := application.NewGetAllAlerts(db)
	getBySensor := application.NewGetBySensorAlert(db)
	ws := application.NewSendAlertUseCase(websocketService)

	// Crear instancia de ProcessSensor
	processSensor := application.NewProcessSensor(rabbitService, save, ws)

	// Crear e inicializar los controladores
	getAllController := handlers.NewGetAllAlerts(getAll)
	getBySensorController := handlers.NewGetBySensor(getBySensor)
	wsHandler := handlers.NewWebSocketHandler(websocketService)

	// Registrar rutas
	SetupRoutes(router, getAllController, getBySensorController, wsHandler)

	// Iniciar WebSocket Server en un goroutine
	go websocketService.Start()
	go func() {
		_, err := rabbitService.FetchReports()
		if err != nil {
			log.Fatalf("Error al obtener reportes de sensores: %v", err)
		}
	}()

	go processSensor.StartProcessingSensors()
}
