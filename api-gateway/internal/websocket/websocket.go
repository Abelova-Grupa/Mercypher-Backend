package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	// "sync"

	"github.com/gorilla/websocket"
)


// Storing data as RawMessage so it can be unmarshaled according to the given type.
type Envelope struct {
    Type string          `json:"type"`
    Data json.RawMessage `json:"data"` // defer decoding of data
}

// Map connections is of structure connection : isActive
// mutex is used because maps aren't thread safe
// type WebSocketConnectionManager struct {
// 	connections map[*websocket.Conn]bool
// 	mu          sync.Mutex
// }

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Accept all origins (for testing).
		return true
	},
}

// func NewWebSocketManager() *WebSocketConnectionManager {
// 	return &WebSocketConnectionManager{
// 		connections: make(map[*websocket.Conn]bool),
// 	}
// }

func Respond(conn *websocket.Conn, messageType int, response string) error {
	if err := conn.WriteMessage(messageType, []byte(response)); err != nil {
		return err	
	} else {
		return nil
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

		// Unmarshal the message
		var env Envelope
		if err := json.Unmarshal(msg, &env); err != nil {
			log.Println("Failed to unmarshall message!")
			continue
		}

		// Get message type and act accordingly
		switch env.Type {
			case "message": 
				if err := Respond(conn, messageType, "Message received!"); err != nil {
					log.Println("Couldn't respond.")
				}
			default:
				if err := Respond(conn, messageType, "Unknown message type!"); err != nil {
					log.Println("Couldn't respond.")
				}
		}

		// Echo it back
		
	}
}
