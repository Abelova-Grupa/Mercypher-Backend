package main

import (
	"log"
	"sync"

	"github.com/Abelova-Grupa/Mercypher/api/internal/servers"

	"github.com/Abelova-Grupa/Mercypher/api/internal/clients"
)

func main() {

	log.Println("Gateway thread started.")

	// wg - A wait group that is keeping the process alive for 3 different routines:
	//		1) Gateway routine
	//		2) gRPC server routine
	//		3) HTTP server routine
	var wg sync.WaitGroup
	wg.Add(3)

	// Servers declaration
	httpServer := servers.NewHttpServer(&wg)
	grpcServer := servers.NewGrpcServer(&wg)

	// Start server routines
	httpServer.Start(":8080")
	grpcServer.Start(":50051")

	// Starting clients to other services.

	// Message client setup
	message_client, err := clients.NewMessageClient("localhost:50052")
	if message_client == nil || err != nil{
		log.Fatalln("Client failed to connect to message service: ", err)
	}
	defer message_client.Close()

	// Relay client setup

	// User client setup
	
	// Session client setup

	// Wait for all routines.
	// Note:	DO NOT PLACE ANY CODE UNDER THE FOLLOWING STATEMENT.
	wg.Wait()
}
