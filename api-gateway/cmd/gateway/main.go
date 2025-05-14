package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Abelova-Grupa/Mercypher/api/internal/handlers"
)

func main() {
	// Da ne uvlacim ceo gin zbog ove jedne rute...
	http.HandleFunc("/ws", handlers.HandleWebSocket)

	fmt.Println("WebSocket server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
