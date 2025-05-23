package websocket

import (
	"encoding/json"
	"log"
	"net/http"

	// "sync"

	"github.com/gorilla/websocket"
)

// Envelope stores data as RawMessage so it can be unmarshaled according to the given type.
type Envelope struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"` // defer decoding of data
}

// ChatMessage stores data of various contents of Envelope.Data json
type ChatMessage struct {
	MessageID  string `json:"message_id"`
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Timestamp  int64  `json:"timestamp"`
	Body       string `json:"body"`
}

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Accept all origins (for testing).
		return true
	},
}

func Respond(conn *websocket.Conn, messageType int, env Envelope) error {

	jsonData, err := json.Marshal(env)

	if err != nil {
		log.Println("Error marshaling message: ", err)
		return err
	}

	if err := conn.WriteMessage(messageType, jsonData); err != nil {
		log.Println("Error writing the response: ", err)
		return err
	}

	return nil

}

func HandleClient(conn *websocket.Conn) {
	defer conn.Close()
	log.Println("New client handler started @", conn.RemoteAddr())

	for {
		// Read a message from the client
		_, msg, err := conn.ReadMessage()

		if err != nil {
			log.Println("Error reading message:", err)
			break
		}


		// Unmarshal the message
		var env Envelope
		if err := json.Unmarshal(msg, &env); err != nil {
			log.Println("Failed to unmarshall message!")
			if err := Respond(conn, websocket.TextMessage, Envelope{"error", nil}); err != nil {
				log.Println("Couldn't respond.")
			}
			continue
		}

		// Get message type and act accordingly
		switch env.Type {
		case "ping":
			if err := Respond(conn, websocket.PongMessage, Envelope{"pong", nil}); err != nil {
				log.Println("Couldn't respond.")
			}
		case "message":
			if err := Respond(conn, websocket.TextMessage, Envelope{"message received", nil}); err != nil {
				log.Println("Couldn't respond.")
			}
		default:
			if err := Respond(conn, websocket.TextMessage, Envelope{"invalid type received", nil}); err != nil {
				log.Println("Couldn't respond.")
			}
		}

	}
}
