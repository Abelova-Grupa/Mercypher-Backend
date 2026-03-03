package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Abelova-Grupa/Mercypher/message-service/internal/config"
	"github.com/Abelova-Grupa/Mercypher/message-service/internal/kafka"
	"github.com/Abelova-Grupa/Mercypher/message-service/internal/repository"
	"github.com/Abelova-Grupa/Mercypher/message-service/internal/server"
	pb "github.com/Abelova-Grupa/Mercypher/proto/message"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	// runing configuration
	config.LoadEnv()
	kafkaBrokerEnv := config.GetEnv("KAFKA_BROKERS", "localhost:9092")
	brokers := strings.Split(kafkaBrokerEnv, ",")
	port := config.GetEnv("PORT", "50052")

	host := config.GetEnv("DB_HOST", "localhost")
	dbPort := config.GetEnv("DB_PORT", "5433")
	user := config.GetEnv("POSTGRES_USER", "mercypher_admin")
	pass := config.GetEnv("POSTGRES_PASSWORD", "password321")
	name := config.GetEnv("POSTGRES_DB", "message_db")
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, dbPort, user, pass, name)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	log.Printf("Connected to postres -> " + dsn)
	defer db.Close()
	repo := repository.NewMessageRepository(db)
	consumer := kafka.NewKafkaConsumer(repo, brokers)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go consumer.Start(ctx)

	// starting a listener
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// starting grpc server with message service
	grpcServer := grpc.NewServer()
	msgServer := server.NewMessageServer(brokers, repo)
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
	consumer.Close()
	log.Println("Server stopped.")
}
