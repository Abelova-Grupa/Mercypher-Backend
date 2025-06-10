package main

import (
	"log"
	"sync"

	"github.com/Abelova-Grupa/Mercypher/api/internal/domain"
	"github.com/Abelova-Grupa/Mercypher/api/internal/servers"
	"github.com/Abelova-Grupa/Mercypher/api/internal/websocket"

	cli "github.com/Abelova-Grupa/Mercypher/api/internal/clients"
)

type Gateway struct {
	// WaitGroup for routine synchronization
	wg				*sync.WaitGroup

	// Websocket registration channels
	register		chan *websocket.Websocket
	unregister		chan *websocket.Websocket
	
	// Channels for communication between Gateway and HTTP/gRPC servers
	inHttp			chan *domain.Envelope
	outHttp			chan *domain.Envelope
	inGrpc			chan *domain.Envelope
	outGrpc			chan *domain.Envelope
	
	// Websocket map for storing connected clients 
	clients     	map[*websocket.Websocket]struct{}
	mu          	sync.RWMutex             

	// Pointers to clients toward other serices
	messageClient	*cli.MessageClient
	relayClient		*cli.RelayClient
	userClient		*cli.UserClient	
	sessionClient	*cli.SessionClient
}

// Gateway Constructor
func NewGateway(wg *sync.WaitGroup, 
	mc *cli.MessageClient, 
	rc *cli.RelayClient, 
	uc *cli.UserClient, 
	sc *cli.SessionClient) *Gateway {
	return &Gateway{
		wg:				wg,
		register: 		make(chan *websocket.Websocket, 32),
		unregister: 	make(chan *websocket.Websocket, 32),
		inHttp:			make(chan *domain.Envelope, 100),
		outHttp:		make(chan *domain.Envelope, 100),
		inGrpc:			make(chan *domain.Envelope, 100),
		outGrpc:		make(chan *domain.Envelope, 100),
		clients: 		make(map[*websocket.Websocket]struct{}),
		messageClient: 	mc,
		relayClient: 	rc,
		userClient: 	uc,
		sessionClient: 	sc,
	}
}

func (g *Gateway) Close() {
	
}

// TODO: Implement gateway message routing here
func (g *Gateway) Start() {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		for {
			select {
			// Handle new websocket connection
			case ws := <-g.register:
				g.mu.Lock()
				g.clients[ws] = struct{}{}
				g.mu.Unlock()
				log.Println("Client registered:", ws.Client.UserId, "\t\t Connected clients: ", len(g.clients))
	
			// Handle websocket disconnection
			case ws := <-g.unregister:
				g.mu.Lock()
				delete(g.clients, ws)
				g.mu.Unlock()
				log.Println("Client unregistered:", ws.Client.UserId, "\t Connected clients: ", len(g.clients))
	
			// Handle HTTP input messages
			case msg := <-g.inHttp:
				log.Println("Received from HTTP:", msg)
				// Add logic to route or process msg
	
			// Handle gRPC input messages
			case msg := <-g.inGrpc:
				log.Println("Received from gRPC:", msg)
				// Add logic to route or process msg
	
			// Handle messages going to HTTP
			case msg := <-g.outHttp:
				log.Println("Sending to HTTP:", msg)
				// Forward to HTTP service
	
			// Handle messages going to gRPC
			case msg := <-g.outGrpc:
				log.Println("Sending to gRPC:", msg)
				// Forward to gRPC service
			}
		}
	}()
}

func main() {
	// wg - A wait group that is keeping the process alive for 3 different routines:
	//		1) Gateway routine
	//		2) gRPC server routine
	//		3) HTTP server routine
	var wg sync.WaitGroup

	// Starting clients to other services.
	// Message client setup
	messageClient, err := cli.NewMessageClient("localhost:50052")
	if messageClient == nil || err != nil{
		log.Fatalln("Client failed to connect to message service: ", err)
	}
	defer messageClient.Close()

	// Relay client setup
	relayClient, err := cli.NewRelayClient("localhost:50053")
	if relayClient == nil || err != nil{
		log.Fatalln("Client failed to connect to relay service: ", err)
	}
	defer relayClient.Close()

	// User client setup
	userClient, err := cli.NewUserClient("localhost:50054")
	if userClient == nil || err != nil{
		log.Fatalln("Client failed to connect to user service: ", err)
	}
	defer userClient.Close()

	// Session client setup
	sessionClient, err := cli.NewSessionClient("localhost:50055")
	if sessionClient == nil || err != nil{
		log.Fatalln("Client failed to connect to session service: ", err)
	}
	defer sessionClient.Close()

	// Servers declaration
	gateway := NewGateway(&wg, messageClient, relayClient, userClient, sessionClient)

	httpServer := servers.NewHttpServer(&wg, gateway.inHttp, gateway.outHttp, gateway.register, gateway.unregister)
	grpcServer := servers.NewGrpcServer(&wg, gateway.inGrpc, gateway.outGrpc)

	// Start server routines
	gateway.Start()

	httpServer.Start(":8080")
	grpcServer.Start(":50051")

	// Wait for all routines.
	// Note:	DO NOT PLACE ANY CODE UNDER THE FOLLOWING STATEMENT.
	wg.Wait()
}
