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

	// Note: 	grpc server has its own weird struct that i don't want to mess with, so
	// 			until i figure out how to make it, there won't be a grpcServer field for
	// 			it will be both created and started in servers.StartGrpcServer method.
	//
	// 			Though I would really like to assign a wait group field to grpc server
	//			instead of passing it as a parameter in start method.

	// Start server routines
	httpServer.Start(":8080")
	servers.StartGrpcServer(":50051", &wg)

	// Starting clients to other services.

	// Message client setup
	message_client, err := clients.NewMessageClient("localhost:50052")
	if message_client == nil || err != nil{
		log.Fatalln("Client failed to connect to message service: ", err)
	}
	defer message_client.Close()

	// Wait for all routines.
	// Note:	DO NOT PLACE ANY CODE UNDER THE FOLLOWING STATEMENT.
	wg.Wait()
}
