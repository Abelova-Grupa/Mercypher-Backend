package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Abelova-Grupa/Mercypher/api-gateway/internal/domain"
	"github.com/Abelova-Grupa/Mercypher/api-gateway/internal/servers"
	"github.com/Abelova-Grupa/Mercypher/api-gateway/internal/websocket"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"

	cli "github.com/Abelova-Grupa/Mercypher/api-gateway/internal/clients"
	cfg "github.com/Abelova-Grupa/Mercypher/api-gateway/internal/config"
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

	// Kafka 
	kafkaIn			chan *domain.ChatMessage
	
	// Websocket map for storing connected clients 
	clients     	map[string]*websocket.Websocket
	mu          	sync.RWMutex             

	// Pointers to clients toward other serices
	messageClient	*cli.MessageClient
	userClient		*cli.UserClient	
	sessionClient	*cli.SessionClient
}

// Gateway Constructor
func NewGateway(wg *sync.WaitGroup, 
	mc *cli.MessageClient, 
	uc *cli.UserClient, 
	sc *cli.SessionClient) *Gateway {
	return &Gateway{
		wg:				wg,
		register: 		make(chan *websocket.Websocket, 32),
		unregister: 	make(chan *websocket.Websocket, 32),
		inHttp:			make(chan *domain.Envelope, 100),
		outHttp:		make(chan *domain.Envelope, 100),
		inGrpc:			make(chan *domain.Envelope, 100),
		kafkaIn: 		make(chan *domain.ChatMessage, 100),
		outGrpc:		make(chan *domain.Envelope, 100),
		clients: 		make(map[string]*websocket.Websocket),
		messageClient: 	mc,
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
				g.clients[ws.Client.UserId] = ws 
				g.mu.Unlock()
				log.Println("Client registered:", ws.Client.UserId, "\t\t Connected clients: ", len(g.clients))

			// Handle websocket disconnection
			case ws := <-g.unregister:
				g.mu.Lock()
				delete(g.clients, ws.Client.UserId)
				g.mu.Unlock()
				log.Println("Client unregistered:", ws.Client.UserId, "\t Connected clients: ", len(g.clients))
				
			case msg := <-g.kafkaIn:
				// TODO: Check if client failed..
				if _, ok := g.clients[msg.Receiver_id]; ok {
					g.clients[msg.Receiver_id].SendChatMessage(*msg)
				}
				if _, ok := g.clients[msg.SenderId]; ok {
					g.clients[msg.SenderId].SendMessageAck(*msg)
				}
			// These might be unnecessary for grpc and http clients can run in separate routines and handle their connections there.

			// // Handle HTTP input messages
			// case msg := <-g.inHttp:
			// 	log.Println("Received from HTTP:", msg)
			// 	// Add logic to route or process msg
	
			// // Handle gRPC input messages
			// case msg := <-g.inGrpc:
			// 	log.Println("Received from gRPC:", msg)
			// 	// Add logic to route or process msg
	
			// // Handle messages going to HTTP
			// case msg := <-g.outHttp:
			// 	log.Println("Sending to HTTP:", msg)
			// 	// Forward to HTTP service
	
			// // Handle messages going to gRPC
			// case msg := <-g.outGrpc:
			// 	log.Println("Sending to gRPC:", msg)
			// 	// Forward to gRPC service

			case env := <-g.inHttp:
                // Handle a message coming FROM a local user to the Gateway
                var chatMsg domain.ChatMessage
                if err := json.Unmarshal(env.Data, &chatMsg); err == nil {
					log.Printf("%s -> %s [ %s ]", chatMsg.SenderId, chatMsg.Receiver_id, chatMsg.Body)
                    chatMsg.Timestamp = time.Now().Unix()
                    g.messageClient.SendMessage(chatMsg)
                } else {
					fmt.Println("Message service failed: ", err)
				}

			}

			// Check channels for each
			// for _, client := range g.clients { 
			// 	select {
			// 	case msg := <-client.Out:
			// 		var chatMsg domain.ChatMessage
			// 		log.Println("Message received at main.")
			// 		if err := json.Unmarshal(msg.Data, &chatMsg); err != nil {
			// 			log.Println("Invalid message format:", err)
			// 			continue
			// 		}

			// 		// Attach message metadata
			// 		chatMsg.SenderId = client.Client.UserId
			// 		chatMsg.Timestamp = time.Now().Unix()

			// 		log.Printf("%s -> %s [ %s ]", chatMsg.SenderId, chatMsg.Receiver_id, chatMsg.Body)

			// 		if err := g.messageClient.SendMessage(chatMsg); err != nil {
			// 			log.Println("Message service failed: ", err)
			// 		}
			// 	}
			// }

		}
	}()
}

func main() {
	// wg - A wait group that is keeping the process alive for 3 different routines:
	//		1) Gateway routine
	//		2) gRPC server routine
	//		3) HTTP server routine
	var wg sync.WaitGroup

	env :=os.Getenv("ENVIRONMENT")
	messageUrl := resolveHost("MESSAGE_HOST","localhost:50052","MESSAGE_PORT",env)
	userUrl := resolveHost("USER_HOST","localhost:50054","USER_PORT",env)
	sessionUrl := resolveHost("SESSION_HOST","localhost:50055","SESSION_PORT",env)
	groupUrl := resolveHost("GROUP_HOST","localhost:50056","GROUP_PORT",env)

	// messageHost := cfg.GetEnv("MESSAGE_HOST", "localhost:50052")
	// userHost := cfg.GetEnv("USER_HOST", "localhost:50054")
	// sessionHost := cfg.GetEnv("SESSION_HOST", "localhost:50055")
	// groupHost := cfg.GetEnv("GROUP_HOST", "localhost:50056")
	kafkaBrokers := cfg.GetEnv("KAFKA_BROKERS", "localhost:9092")
	httpPort := cfg.GetEnv("HTTP_PORT", "8080")
	grpcPort := cfg.GetEnv("GRPC_PORT", "50051")

	log.Printf("========== Environment ==========")
	log.Printf("Message host:\t\t %v", messageUrl)
	log.Printf("User host:\t\t %v", userUrl)
	log.Printf("Session host:\t\t %v", sessionUrl)
	// log.Printf("Kafka brokers:\t\t %v\n", kafkaBrokers)

	// Starting clients to other services.
	// Message client setup
	messageClient, err := cli.NewMessageClient(messageUrl)
	if messageClient == nil || err != nil{
		log.Fatalln("Client failed to connect to message service: ", err)
	}
	defer messageClient.Close()

	// User client setup
	userClient, err := cli.NewUserClient(userUrl)
	if userClient == nil || err != nil {
		log.Fatalln("Client failed to connect to user service: ", err)
	}
	defer userClient.Close()

	// Session client setup
	sessionClient, err := cli.NewSessionClient(sessionUrl)
	if sessionClient == nil || err != nil {
		log.Fatalln("Client failed to connect to session service: ", err)
	}
	defer sessionClient.Close()

	groupClient, err := cli.NewGroupClient(groupUrl)
	if groupClient == nil || err != nil {
		log.Fatalln("Client failed to connect to group service: ", err)
	}

	// Servers declaration
	gateway := NewGateway(&wg, messageClient, userClient, sessionClient)

	if env == "azure" {
		namespace := os.Getenv("AZURE_SERVICE_BUS_CONN_STR")
		if namespace == "" {
			panic("azure service bus connection string not loaded")
		}
		clientOpt := &azservicebus.ClientOptions{
			TLSConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
			RetryOptions: azservicebus.RetryOptions{
				MaxRetries: 5,
			},
		}
		client, err := azservicebus.NewClientFromConnectionString(namespace, clientOpt)
		if err != nil {
			panic("unable to create azure message server")
		}
		
		busReceiver, err := servers.CreateTopicReceiver(client,"azure-topic","gateway-sub",nil)
		if err != nil {
			panic("couldn't create a azure busReceiver struct")
		}

		azureConsumer := servers.NewBusReceiver(busReceiver,gateway.kafkaIn)
		ctx := context.Background()
		go azureConsumer.Start(ctx)
		defer azureConsumer.Close(ctx)

	}else {
		brokers := strings.Split(kafkaBrokers, ",")
		kafka := servers.NewKafkaConsumer(brokers, "chat-messages-v1", "gw-consumer", gateway.kafkaIn)
		go kafka.StartLiveForwarder(context.Background())
		defer kafka.Close()
	}

	httpServer := servers.NewHttpServer(&wg, gateway.inHttp, 
		gateway.outHttp, 
		gateway.register, 
		gateway.unregister, 
		userClient, 
		sessionClient,
		messageClient,
		groupClient)
	grpcServer := servers.NewGrpcServer(&wg, gateway.inGrpc, gateway.outGrpc)

	// Start server routines
	gateway.Start()

	httpServer.Start("0.0.0.0:"+httpPort)

	

	grpcServer.Start("0.0.0.0:"+grpcPort)

	// Wait for all routines.
	// Note:	DO NOT PLACE ANY CODE UNDER THE FOLLOWING STATEMENT.
	wg.Wait()
}

func resolveHost(hostkey string, defaultHost string,portKey string, env string) string{
	if os.Getenv(hostkey) == "" {
		return defaultHost + ":" + os.Getenv(portKey)
	}
	//Internal connection to grpc services
	if env == "azure" {
		return fmt.Sprintf("%s:%s",os.Getenv(hostkey),os.Getenv(portKey))
	}

	return os.Getenv(hostkey)
}
