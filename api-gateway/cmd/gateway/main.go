package main

import (
	"log"

	server "github.com/Abelova-Grupa/Mercypher/api/internal/server"
)

func main() {
	server := server.InitServer()
	log.Fatal(server.Start(":8080"))
}
