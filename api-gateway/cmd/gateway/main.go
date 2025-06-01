package main

import (
	"log"
	"sync"

	"github.com/Abelova-Grupa/Mercypher/api/internal/domain"
	"github.com/Abelova-Grupa/Mercypher/api/internal/servers"
	"github.com/Abelova-Grupa/Mercypher/api/internal/websocket"

	"github.com/Abelova-Grupa/Mercypher/api/internal/clients"
)

type Gateway struct {
	wg			*sync.WaitGroup
	register	chan *websocket.Websocket
	unregister	chan *websocket.Websocket
	inHttp		chan *domain.Envelope
	outHttp		chan *domain.Envelope
	inGrpc		chan *domain.Envelope
	outGrpc		chan *domain.Envelope
}

func NewGateway(wg *sync.WaitGroup) *Gateway {
	return &Gateway{
		wg:				wg,
		register: 		make(chan *websocket.Websocket, 32),
		unregister: 	make(chan *websocket.Websocket, 32),
		inHttp:			make(chan *domain.Envelope, 100),
		outHttp:		make(chan *domain.Envelope, 100),
		inGrpc:			make(chan *domain.Envelope, 100),
		outGrpc:		make(chan *domain.Envelope, 100),

	}
}

func (g *Gateway) Start() {
	g.wg.Add(1)
	go func(){
		defer g.wg.Done()
		for msg := range g.inGrpc {
			log.Println("Gateway received channel message:", msg.Type, msg.Data)
		}
	}()
}

func main() {


	log.Println("Gateway thread started.")

	// wg - A wait group that is keeping the process alive for 3 different routines:
	//		1) Gateway routine
	//		2) gRPC server routine
	//		3) HTTP server routine
	var wg sync.WaitGroup

	// Servers declaration
	gateway := NewGateway(&wg)

	httpServer := servers.NewHttpServer(&wg, gateway.inHttp, gateway.outHttp)
	grpcServer := servers.NewGrpcServer(&wg, gateway.inGrpc, gateway.outGrpc)

	// Start server routines
	gateway.Start()

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

	// Handle system traffic
	
	// Wait for all routines.
	// Note:	DO NOT PLACE ANY CODE UNDER THE FOLLOWING STATEMENT.
	wg.Wait()
}
