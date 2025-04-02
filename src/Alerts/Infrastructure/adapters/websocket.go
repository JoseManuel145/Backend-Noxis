package adapters

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketAdapter struct {
	clients  map[string]map[*websocket.Conn]bool // Mapa de sensores a clientes conectados
	mu       sync.Mutex
	upgrader websocket.Upgrader
}

func NewWebSocketAdapter() *WebSocketAdapter {
	return &WebSocketAdapter{
		clients: make(map[string]map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

// Manejar conexiones WebSocket para sensores específicos
func (ws *WebSocketAdapter) HandleConnections(sensor string, w http.ResponseWriter, r *http.Request) {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error al actualizar WebSocket:", err)
		return
	}

	ws.mu.Lock()
	if ws.clients[sensor] == nil {
		ws.clients[sensor] = make(map[*websocket.Conn]bool)
	}
	ws.clients[sensor][conn] = true
	ws.mu.Unlock()

	log.Printf("Nuevo cliente suscrito al sensor %s", sensor)

	defer func() {
		ws.mu.Lock()
		delete(ws.clients[sensor], conn)
		ws.mu.Unlock()
		conn.Close()
		log.Printf("Cliente desconectado de %s", sensor)
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// Enviar mensaje solo a los clientes del sensor correspondiente
func (ws *WebSocketAdapter) SendMessage(sensor string, message []byte) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if clients, ok := ws.clients[sensor]; ok {
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("❌ Error enviando mensaje WebSocket:", err)
				client.Close()
				delete(clients, client)
			} else {
				log.Printf("📡 Mensaje enviado a %s: %s", sensor, string(message))
			}
		}
	}
}

func (ws *WebSocketAdapter) Start() {
	log.Println("Servidor WebSocket en marcha...")
	http.HandleFunc("/ws/", func(w http.ResponseWriter, r *http.Request) {
		// Extraer el sensor de la URL
		sensor := r.URL.Path[len("/ws/"):]
		if sensor == "" {
			http.Error(w, "Sensor no especificado", http.StatusBadRequest)
			return
		}
		ws.HandleConnections(sensor, w, r)
	})

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Println("Error iniciando WebSocket:", err)
	}
}
