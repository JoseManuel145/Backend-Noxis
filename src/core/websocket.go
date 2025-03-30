package core

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebSocketServer maneja las conexiones WebSocket
type WebSocketServer struct {
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
	upgrader  websocket.Upgrader
}

// NewWebSocketServer inicializa WebSocketServer
func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

// HandleConnections maneja las conexiones WebSocket
func (ws *WebSocketServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error al actualizar WebSocket:", err)
		return
	}
	defer conn.Close()
	ws.clients[conn] = true

	// Escuchar mensajes (opcional)
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			delete(ws.clients, conn)
			break
		}
	}
}

// Start inicia el servidor WebSocket
func (ws *WebSocketServer) Start() {
	http.HandleFunc("/ws", ws.HandleConnections)
	go func() {
		if err := http.ListenAndServe(":8081", nil); err != nil {
			fmt.Println("Error iniciando WebSocket:", err)
		}
	}()
	fmt.Println("WebSocket corriendo en ws://localhost:8081/ws")
}
