package main

import (
	//"context"
	"fmt"
	"log"
	"net"

	//"time"

	"github.com/Abelova-Grupa/Mercypher/message-service/internal/config"
	messagepb "github.com/Abelova-Grupa/Mercypher/message-service/internal/grpc"
	"github.com/Abelova-Grupa/Mercypher/message-service/internal/model"
	"github.com/Abelova-Grupa/Mercypher/message-service/internal/repository"
	"github.com/Abelova-Grupa/Mercypher/message-service/internal/server"
	"github.com/Abelova-Grupa/Mercypher/message-service/internal/service"

	"google.golang.org/grpc"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getDatabaseParameters() string {
	config.LoadEnv()

	user := config.GetEnv("DB_USER", "postgres")
	pass := config.GetEnv("DB_PASSWORD", "")
	host := config.GetEnv("DB_HOST", "127.0.0.1")
	port := config.GetEnv("DB_PORT", "5432")
	name := config.GetEnv("DB_NAME", "mercypher_msg")
	ssl := config.GetEnv("DB_SSLMODE", "disable")
	tz := config.GetEnv("DB_TIMEZONE", "UTC")

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&timezone=%s",
		user, pass, host, port, name, ssl, tz,
	)
}

// TODO: Move to config?
func connect(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	log.Println("Attempting to connect to the messages database...")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	} else {
		log.Println("Connected to the users database.")
	}

	// Auto-migrate
	if err := db.AutoMigrate(&model.ChatMessage{}); err != nil {
		log.Fatal("auto-migration failed:", err)
	}

	return db
}

func main() {
	conn := connect(getDatabaseParameters())
	repo := repository.NewMessageRepository(conn)
	service := service.NewMessageService(repo)
	server := server.NewMessageServer(service)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	messagepb.RegisterMessageServiceServer(grpcServer, server)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	// server := server.NewMessageServer()
	// var msg model.ChatMessage
	// msg.Message_id = "test_DZJLKAFSKJGDKJHFGJHLI"
	// msg.Sender_id = "testuser2"
	// msg.Receiver_id = "testuser1"
	// msg.Body = "Proba proba jen dva tri"
	// msg.Timestamp = time.Now()
	// repo.CreateMessage(context.Background(), &msg)

	// ######################
	// 	Code used for texting relay client
	// ctx := context.Background()
	// channel := make(chan func(), 5)
	// defer close(channel)
	// go relay_client.StartClient(channel)

	// scanner := bufio.NewScanner(os.Stdin)
	// for {
	// 	fmt.Print("Enter a message (or 'exit' to quit): ")
	// 	if !scanner.Scan() {
	// 		break
	// 	}
	// 	text := scanner.Text()

	// 	if text == "exit" {
	// 		break
	// 	}
	// 	msg := model.ChatMessage{
	// 		Message_id:  "123",
	// 		Sender_id:   "steven",
	// 		Receiver_id: "derek",
	// 		Body:        text,
	// 		Timestamp:   time.Now(),
	// 	}

	// 	// Send a function that prints the message
	// 	channel <- func() {
	// 		_ = relay_client.RelayMessage(ctx, &msg)
	// 	}
	// }
}
