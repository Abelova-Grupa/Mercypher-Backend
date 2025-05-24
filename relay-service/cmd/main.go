package main

import (
	"github.com/Abelova-Grupa/Mercypher-Backend/relay-service/internal/config"
	"github.com/Abelova-Grupa/Mercypher-Backend/relay-service/internal/server"
)

func main() {
	// Loading env file
	config.LoadEnv()

	// Starting server
	server.StartServer()
}
