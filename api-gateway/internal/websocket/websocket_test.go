package websocket

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Abelova-Grupa/Mercypher/api-gateway/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func setupTestWebsocket(t *testing.T) (*websocket.Conn, *Websocket, chan *domain.Envelope, chan *Websocket) {
	gin.SetMode(gin.TestMode)

	in := make(chan *domain.Envelope, 1)
	unregister := make(chan *Websocket, 1)

	// Create HTTP test server that upgrades to websocket
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := Upgrader.Upgrade(w, r, nil)
		assert.NoError(t, err)

		ws := NewWebsocket(conn, domain.User{
			UserId:   "alice",
			Username: "alice",
		}, unregister, in)

		go ws.HandleClient()
	}))

	// Convert http://127.0.0.1 → ws://127.0.0.1
	url := "ws" + server.URL[len("http"):]

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.NoError(t, err)

	return conn, nil, in, unregister
}

func TestWebsocket_MessageForwarding(t *testing.T) {
	conn, _, in, _ := setupTestWebsocket(t)
	defer conn.Close()

	chat := domain.ChatMessage{
		Body: "Hello",
	}

	msgData, _ := json.Marshal(chat)

	env := domain.Envelope{
		Type: "message",
		Data: msgData,
	}

	data, _ := json.Marshal(env)

	err := conn.WriteMessage(websocket.TextMessage, data)
	assert.NoError(t, err)

	select {
	case received := <-in:
		assert.Equal(t, "message", received.Type)
	case <-time.After(time.Second):
		t.Fatal("Message was not forwarded to In channel")
	}
}

func TestWebsocket_InvalidJSON(t *testing.T) {
	conn, _, _, _ := setupTestWebsocket(t)
	defer conn.Close()

	err := conn.WriteMessage(websocket.TextMessage, []byte("invalid-json"))
	assert.NoError(t, err)

	_, resp, err := conn.ReadMessage()
	assert.NoError(t, err)

	var result domain.Envelope
	json.Unmarshal(resp, &result)

	assert.Equal(t, "error", result.Type)
}