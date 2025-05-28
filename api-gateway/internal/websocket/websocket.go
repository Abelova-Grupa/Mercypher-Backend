package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/Abelova-Grupa/Mercypher/api/internal/domain"
	// "sync"

	"github.com/gorilla/websocket"
)


var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Accept all origins (for testing).
		return true
	},
}

func Respond(conn *websocket.Conn, messageType int, env domain.Envelope) error {

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
		var env domain.Envelope
		if err := json.Unmarshal(msg, &env); err != nil {
			log.Println("Failed to unmarshall message!")
			if err := Respond(conn, websocket.TextMessage, domain.Envelope{Type: "error", Data: nil}); err != nil {
				log.Println("Couldn't respond.")
			}
			continue
		}

		// Get message type and act accordingly
		switch env.Type {
		case "ping":
			if err := Respond(conn, websocket.PongMessage, domain.Envelope{Type: "pong", Data: nil}); err != nil {
				log.Println("Couldn't respond.")
			}
		case "message":
			if err := Respond(conn, websocket.TextMessage, domain.Envelope{Type: "message received", Data: nil}); err != nil {
				log.Println("Couldn't respond.")
			}
		default:
			if err := Respond(conn, websocket.TextMessage, domain.Envelope{Type: "invalid type received", Data: nil}); err != nil {
				log.Println("Couldn't respond.")
			}
		}

	}
}
