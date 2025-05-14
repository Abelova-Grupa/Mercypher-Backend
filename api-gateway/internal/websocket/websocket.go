package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Map connections is of structure connection : isActive
// mutex is used because maps aren't thread safe
type WebSocketConnectionManager struct {
	connections map[*websocket.Conn]bool
	mu          sync.Mutex
}

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Accept all origins (for testing).
		return true
	},
}

func NewWebSocketManager() *WebSocketConnectionManager {
	return &WebSocketConnectionManager{
		connections: make(map[*websocket.Conn]bool),
	}
}

func HandleClient(conn *websocket.Conn) {
	defer conn.Close()
	log.Println("New client handler started @", conn.RemoteAddr())

	for {
		// Read a message from the client
		messageType, msg, err := conn.ReadMessage()

		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		log.Printf("Received: %s", msg)

		// Echo it back
		err = conn.WriteMessage(messageType, msg)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}
