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

	// Registrar rutas
	SetupRoutes(router, getAllController, getBySensorController)

	// Iniciar WebSocket Server en una goroutine
	go func() {
		log.Println("Iniciando servidor WebSocket...")
		websocketService.Start()
	}()

	// Iniciar procesamiento de sensores en una goroutine
	go func() {
		log.Println("Iniciando procesamiento de sensores...")
		processSensor.StartProcessingSensors()
	}()

	// Mantener el flujo principal activo
	select {}
}
