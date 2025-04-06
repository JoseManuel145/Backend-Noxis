package adapters

import (
	"Backend/src/Alerts/domain"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketAdapter struct {
	connections map[*websocket.Conn]bool // Mapa de conexiones activas
	mu          sync.Mutex
	upgrader    websocket.Upgrader
}

func NewWebSocketAdapter() *WebSocketAdapter {
	return &WebSocketAdapter{
		connections: make(map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (ws *WebSocketAdapter) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("❌ Error al actualizar protocolo:", err)
		return
	}
	defer conn.Close()

	// Registrar la conexión
	ws.mu.Lock()
	ws.connections[conn] = true
	ws.mu.Unlock()
	log.Println("✅ Nueva conexión registrada")

	// Manejar mensajes entrantes
	for {
		var alert domain.Alert
		err := conn.ReadJSON(&alert)
		if err != nil {
			log.Println("❌ Error al leer mensaje JSON:", err)
			break
		}
		log.Printf("⚠️ Alerta recibida -> Sensor: %s, Datos: %+v\n", alert.Sensor, alert.Data)
	}

	// Eliminar la conexión al cerrarse
	ws.mu.Lock()
	delete(ws.connections, conn)
	ws.mu.Unlock()
	log.Println("❌ Conexión cerrada")
}

func (ws *WebSocketAdapter) SendMessage(alert *domain.Alert) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	for client := range ws.connections {
		err := client.WriteJSON(alert)
		if err != nil {
			log.Printf("❌ Error enviando mensaje: %v", err)
			client.Close()
			delete(ws.connections, client)
		} else {
			log.Printf("📤 Mensaje enviado a cliente: %+v", alert)
		}
	}
	return nil
}

func (ws *WebSocketAdapter) Start() {
	log.Println("Servidor WebSocket en marcha...")
	http.HandleFunc("/ws", ws.HandleConnections)

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Println("Error iniciando WebSocket:", err)
	}
}
