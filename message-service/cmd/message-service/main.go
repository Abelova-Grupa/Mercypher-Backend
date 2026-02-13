package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Abelova-Grupa/Mercypher/message-service/internal/config"
	"github.com/Abelova-Grupa/Mercypher/message-service/internal/server"
	pb "github.com/Abelova-Grupa/Mercypher/proto/message"
	"google.golang.org/grpc"
)

func main() {
	// runing configuration
	config.LoadEnv()
	kafkaBrokerEnv := config.GetEnv("KAFKA_BROKERS", "localhost:9092")
	brokers := strings.Split(kafkaBrokerEnv, ",")
	port := config.GetEnv("PORT", "50052")

	// starting a listener
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// starting grpc server with message service
	grpcServer := grpc.NewServer()
	msgServer := server.NewMessageServer(brokers)
	pb.RegisterMessageServiceServer(grpcServer, msgServer)

	// running a server in a goroutine as so graceful shutdown is possible (gemini go brr)
	go func() {
		log.Printf("Message Service is running on port %s...", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Starting kafka consumer for live message forwarding messages
	// gatewayAdr := config.GetEnv("GATEWAY_ADDRESS", "localhost:50051") // if set then its running in a container, otherwise locally
	// go kafka.StartLiveForwarder(context.Background(), brokers, gatewayAdr)

	// Graceful Shutdown (gemini go brr)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("Server stopped.")
}
