package adapters

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketAdapter struct {
	clients  map[*websocket.Conn]bool
	mu       sync.Mutex
	upgrader websocket.Upgrader
}

func NewWebSocketAdapter() *WebSocketAdapter {
	return &WebSocketAdapter{
		clients: make(map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (ws *WebSocketAdapter) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error al actualizar WebSocket:", err)
		return
	}
	defer conn.Close()

	ws.mu.Lock()
	ws.clients[conn] = true
	ws.mu.Unlock()
	log.Println("Nuevo cliente conectado")

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			ws.mu.Lock()
			delete(ws.clients, conn)
			ws.mu.Unlock()
			log.Println("Cliente desconectado")
			break
		}

		log.Printf("Mensaje recibido: %s", msg)
		ws.SendMessage(msg)
	}
}

func (ws *WebSocketAdapter) SendMessage(message []byte) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	for client := range ws.clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Error enviando mensaje WebSocket:", err)
			client.Close()
			delete(ws.clients, client)
		}
	}
}

func (ws *WebSocketAdapter) Start() {
	http.HandleFunc("/ws", ws.HandleConnections)
	go func() {
		if err := http.ListenAndServe(":8081", nil); err != nil {
			log.Println("Error iniciando WebSocket:", err)
		}
	}()
	fmt.Println("WebSocket corriendo en ws://localhost:8081/ws")
}
