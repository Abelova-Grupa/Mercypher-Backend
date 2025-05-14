package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Accept all origins (for testing).
		return true
	},
}

// Simple handler which echoes the message back to the client
func handleClient(conn *websocket.Conn) {
	defer conn.Close()
	log.Println("New client handler started @", conn.RemoteAddr())

	for {
		// Read a message from the client
		messageType, msg, err := conn.ReadMessage()

		if err != nil {
			log.Println("Client disconnected:", err)
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

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	// Handle this client in a new goroutine
	go handleClient(conn)
}

func main() {
	// Da ne uvlacim ceo gin zbog ove jedne rute...
	http.HandleFunc("/ws", handleWebSocket)

	fmt.Println("WebSocket server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
