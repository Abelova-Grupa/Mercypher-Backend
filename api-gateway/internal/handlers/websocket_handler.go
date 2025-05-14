package handlers

import (
	"log"
	"net/http"

	"github.com/Abelova-Grupa/Mercypher/api/internal/websocket"
)

// Simple handler which echoes the message back to the client

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := websocket.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	// Handle this client in a new goroutine
	go websocket.HandleClient(conn)
}
