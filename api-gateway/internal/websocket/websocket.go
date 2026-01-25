package websocket

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Abelova-Grupa/Mercypher/api-gateway/internal/domain"

	// "sync"

	"github.com/gorilla/websocket"
)

//Websocket that serves a logged user.
type Websocket struct {
	Conn 		*websocket.Conn
	Client 		domain.User
	In			chan *domain.Envelope
	Out			chan *domain.Envelope
	unregister	chan *Websocket
}

func NewWebsocket(conn *websocket.Conn, client domain.User, unregister chan *Websocket, in chan *domain.Envelope) *Websocket {
	return &Websocket{
		Conn: 	conn,
		Client: client,
		In:		in,
		Out: 	make(chan *domain.Envelope, 100),
		unregister: unregister,
	}
}

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Accept all origins (for testing).
		return true
	},
}

func (s *Websocket) Respond(messageType int, env domain.Envelope) error {

	jsonData, err := json.Marshal(env)

	if err != nil {
		log.Println("Error marshaling message: ", err)
		return err
	}

	if err := s.Conn.WriteMessage(messageType, jsonData); err != nil {
		log.Println("Error writing the response: ", err)
		return err
	}

	return nil

}

// So, it appears that golang doesn't support function overloading...
func (s *Websocket) SendChatMessage(msg domain.ChatMessage) error {
	
	jsonMessage, err := json.Marshal(msg)

	if err != nil {
		log.Println("Error marshaling message: ", err)
		return err
	}

	env := domain.Envelope{Type: "message", Data: jsonMessage}

	return s.Respond(websocket.TextMessage, env)

}

func (s *Websocket) HandleClient() {
	defer s.Conn.Close()
	log.Println("New client handler started @", s.Conn.RemoteAddr())

	for {
		// Read a message from the client
		_, msg, err := s.Conn.ReadMessage()

		if err != nil {
			// Check whether the user has disconnected from websocket
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
                s.unregister <- s
				break
            } else {
				log.Println("Error reading message:", err)
			}
		}

		// Unmarshal the message
		var env domain.Envelope
		if err := json.Unmarshal(msg, &env); err != nil {
			log.Println("Failed to unmarshall message!")
			if err := s.Respond(websocket.TextMessage, domain.Envelope{Type: "error", Data: nil}); err != nil {
				log.Println("Couldn't respond.")
			}
			continue
		}

		// Get message type and act accordingly
		switch env.Type {
		case "ping":
			if err := s.Respond(websocket.PongMessage, domain.Envelope{Type: "pong", Data: nil}); err != nil {
				log.Println("Couldn't respond.")
			}
		case "message": {
			log.Println("Message received on ws.")
			if err := s.Respond(websocket.TextMessage, domain.Envelope{Type: "message received", Data: nil}); err != nil {
				log.Println("Couldn't respond.")
			}

			// Deconstruct, add metedata, package again.
			// Temporary solution
			var chatMsg domain.ChatMessage
					if err := json.Unmarshal(env.Data, &chatMsg); err != nil {
						log.Println("Invalid message format:", err)
						continue
					}
			chatMsg.SenderId = s.Client.UserId
			env.Data, _ = json.Marshal(chatMsg)
			// FIX NAMING OF THESE
			s.In <- &env
		}
		default:
			if err := s.Respond(websocket.TextMessage, domain.Envelope{Type: "invalid type received", Data: nil}); err != nil {
				log.Println("Couldn't respond.")
			}
		}

	}
}
