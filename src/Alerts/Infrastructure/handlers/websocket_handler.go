package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Estructura del handler WebSocket
type WebSocketHandler struct {
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]bool
}

// Constructor para inicializar el handler
func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		clients: make(map[*websocket.Conn]bool),
	}
}

// Método para manejar conexiones WebSocket
func (wsh *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	conn, err := wsh.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error al establecer conexión WebSocket: %v", err)
		return
	}
	defer conn.Close()

	// Agregar cliente a la lista
	wsh.clients[conn] = true
	log.Println("Nuevo cliente conectado")

	// Leer mensajes entrantes
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error al leer mensaje:", err)
			delete(wsh.clients, conn)
			break
		}

		log.Printf("Mensaje recibido: %s", msg)

		// Reenviar mensaje a todos los clientes conectados
		for client := range wsh.clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("Error al enviar mensaje:", err)
				client.Close()
				delete(wsh.clients, client)
			}
		}
	}
}
