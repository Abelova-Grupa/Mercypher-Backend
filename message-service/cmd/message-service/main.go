package main

import (
	"context"
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
	"github.com/Abelova-Grupa/Mercypher/message-service/internal/store"
	pb "github.com/Abelova-Grupa/Mercypher/proto/message"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	// runing configuration
	config.LoadEnv()
	port := config.GetEnv("PORT", "50052")

	// host := config.GetEnv("DB_HOST", "localhost")
	// dbPort := config.GetEnv("DB_PORT", "5433")
	// user := config.GetEnv("POSTGRES_USER", "mercypher_admin")
	// pass := config.GetEnv("POSTGRES_PASSWORD", "password321")
	// name := config.GetEnv("POSTGRES_DB", "message_db")
	// dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	// 	host, dbPort, user, pass, name)
	// db, err :=
	// db, err := sql.Open("postgres", dsn)
	// if err != nil {
	// 	log.Fatalf("Failed to connect to DB: %v", err)
	// }
	// defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Switched from sql to gorm
	db, err := store.NewMessageDB(ctx)
	if err != nil {
		panic(err)
	}

	dbConn, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer dbConn.Close()
	// log.Printf("Connected to postres -> " + dsn)
	repo := repository.NewMessageRepository(db)

	var msgServer pb.MessageServiceServer
	if os.Getenv("ENVIRONMENT") == "azure" {
		msgServer = server.NewAzureMessageServer(ctx, repo)
		azureServer, ok := msgServer.(*server.AzureMessageServer)
		if !ok {
			panic("cannot cast message server to azure server in azure environment")
		}
		go azureServer.QueueReceiver.Start(ctx)
		defer azureServer.QueueReceiver.Close(ctx)

	} else {
		kafkaBrokerEnv := config.GetEnv("KAFKA_BROKERS", "localhost:9092")
		brokers := strings.Split(kafkaBrokerEnv, ",")
		consumer := kafka.NewKafkaConsumer(repo, brokers)
		go consumer.Start(ctx)
		defer consumer.Close()

		msgServer = server.NewKafkaMessageServer(brokers, repo)
	}

	// starting a listener
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// starting grpc server with message service
	grpcServer := grpc.NewServer()
	pb.RegisterMessageServiceServer(grpcServer, msgServer)

	go func() {
		log.Printf("Message Service is running on port %s...", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
			cancel()
		}
	}()

	// Starting kafka consumer for live message forwarding messages
	// gatewayAdr := config.GetEnv("GATEWAY_ADDRESS", "localhost:50051") // if set then its running in a container, otherwise locally
	// go kafka.StartLiveForwarder(context.Background(), brokers, gatewayAdr)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("Server stopped.")
}
