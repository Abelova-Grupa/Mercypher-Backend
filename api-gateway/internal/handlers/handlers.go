package handlers

import (
	"log"
	"net/http"

	"github.com/Abelova-Grupa/Mercypher/api/internal/websocket"
	"github.com/gin-gonic/gin"
)

func HandleLogin(ctx *gin.Context) {
	log.Println("Not implemented yet!")
	// TODO: Implement
}

func HandleLogout(ctx *gin.Context) {
	log.Println("Not implemented yet!")
	// TODO: Implement
}

func HandleRegister(ctx *gin.Context) {
	log.Println("Not implemented yet!")
	// TODO: Implement
}

func HandleSearchUser(ctx *gin.Context) {
	log.Println("Not implemented yet!")

	username := ctx.Query("username")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error":"username query param is required"})
	}
	// TODO: Implement
}

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
