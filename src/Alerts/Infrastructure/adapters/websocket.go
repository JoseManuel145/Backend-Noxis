package adapters

import (
	"fmt"
	"net/http"

	"Backend/src/Alerts/domain/repositories"

	"github.com/gorilla/websocket"
)

type WebSocketAdapter struct {
	clients  map[*websocket.Conn]bool
	upgrader websocket.Upgrader
}

func NewWebSocketAdapter() repositories.IWebSocketRepository {
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
		fmt.Println("Error al actualizar WebSocket:", err)
		return
	}
	defer conn.Close()
	ws.clients[conn] = true

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			delete(ws.clients, conn)
			break
		}
	}
}

func (ws *WebSocketAdapter) SendMessage(message []byte) {
	for client := range ws.clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			fmt.Println("Error enviando mensaje WebSocket:", err)
			client.Close()
			delete(ws.clients, client)
		}
	}
}

func (ws *WebSocketAdapter) Start() {
	http.HandleFunc("/ws", ws.HandleConnections)
	go func() {
		if err := http.ListenAndServe(":8081", nil); err != nil {
			fmt.Println("Error iniciando WebSocket:", err)
		}
	}()
	fmt.Println("WebSocket corriendo en ws://localhost:8081/ws")
}
