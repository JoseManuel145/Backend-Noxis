package handlers

import (
	"Backend/src/Alerts/Infrastructure/adapters"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebSocketHandler struct {
	wsAdapter *adapters.WebSocketAdapter
}

// Constructor
func NewWebSocketHandler(wsAdapter *adapters.WebSocketAdapter) *WebSocketHandler {
	return &WebSocketHandler{wsAdapter: wsAdapter}
}

// Manejar conexión WebSocket específica por sensor
func (wsh *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	sensor := c.Param("sensor")
	if sensor == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sensor no especificado"})
		return
	}

	wsh.wsAdapter.HandleConnections(sensor, c.Writer, c.Request)
	log.Printf("Cliente conectado al WebSocket de %s", sensor)
}
