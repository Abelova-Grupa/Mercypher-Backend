package handlers

import (
	"log"

	"github.com/Abelova-Grupa/Mercypher/api/internal/websocket"
	"github.com/gin-gonic/gin"
)

// Simple handler which echoes the message back to the client

func HandleWebSocket(ctx *gin.Context) {
	// Upgrade HTTP connection to WebSocket
	conn, err := websocket.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	// Handle this client in a new goroutine
	go websocket.HandleClient(conn)
}
